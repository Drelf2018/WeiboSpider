import os
from datetime import datetime
from io import BytesIO
from math import ceil

import qrcode
import requests
from lxml.etree import HTML
from PIL import Image, ImageDraw, ImageFont

headers = {
    'Connection': 'keep-alive',
    'Accept-Language': 'zh-CN,zh;q=0.9',
    'Accept-Encoding': 'gzip, deflate, br',
    'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8',
    'Upgrade-Insecure-Requests': '1',
    'User-Agent': 'Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/62.0.3202.62 Safari/537.36',
    'cookie': ''
}

def get_content_text(text: str) -> str:
        text = text.replace('<span class="surl-text">', '').replace('</span>', '')
        span = HTML('<span class="ctt">' + text + '</span>')
        for _img in span.xpath('./span[@class="url-icon"]/img'):
            alt, src = _img.xpath('./@alt')[0], _img.xpath('./@src')[0]
            text = text.replace(
                f'<span class="url-icon"><img alt="{alt}" src="{src}" style="width:1em; height:1em;" /></span>',
                alt
            )
        for _a in span.xpath('.//a'):
            href = _a.xpath('./@href')[0]
            atext = _a.xpath('./text()')[0]
            text = text.replace(f'<a href="{href.replace("&", "&amp;")}">{atext}</a>', atext)
            text = text.replace(f'<a  href="{href}" data-hide="">{atext}</a>', atext)
        text = text.replace('<br />', '\n')
        dot = len(text)
        for i in range(dot, 0, -1):
            if not text[i-1] == ' ':
                dot = i
                break
        return text[:dot]

def create_new_img(post: dict, userInfo: dict, cookie: str = None, w: int = 1080) -> Image.Image.save:
    """
    根据微博博文生成图片
    Args:
        post (dict) : 博文信息
        userInfo(dict) : 博主信息
        w (int) : 生成图片的宽度
    Returns:
        PIL 类中的 Image 对象
    """
    headers['cookie'] = cookie

    # 定义字体
    font_type = 'gothon/HarmonyOS_Sans_SC_Regular.ttf'  # 'C:/Windows/Fonts/msyh.ttc'
    font_blod_type = 'C:/Windows/Fonts/msyhbd.ttc'
    font_song = 'C:/Windows/Fonts/simsun.ttc'
    name_font = ImageFont.truetype(font_blod_type, 100)
    text_font = ImageFont.truetype(font_type, 100)

    text_width = 0.8  # 一个系数之后会用到
    # 人为规定一行写 21 个字好看，测量这 21 个字的长度
    # 再根据设定图片宽度 w 乘预设系数 text_width 调整字号
    tw, th = text_font.getsize('我不动脑子随手一写就是标标准准的二十一个字')
    text_font = ImageFont.truetype(font_type, int(100*text_width*w/tw))

    # 将整段文本分成一行行
    split_num = 0
    now_str = ''
    now_height = [[0]]  # 列表中包含多个梓列表，子列表表示这行字所处高度(像素值)以及这行文本内容
    if post.get('repo'):
        try:
            post['repo'] = get_content_text(post['repo'])
        except Exception as e:
            print('repo Error:', e)
        for chn in post['repo']:
            # 每次新增一个文字 如果超过线宽 就把之前的文字存起来
            if chn == '\n' or text_font.getsize(now_str+chn)[0] > text_width*w:
                now_height[-1].append(now_str)
                if not chn == '\n':
                    now_str = chn
                else:
                    now_str = ''
                now_height.append([now_height[-1][0] + text_font.getsize(now_str+chn)[1] + 10])
            else:
                now_str += chn
        now_height[-1].append(now_str)
        now_height.append([now_height[-1][0] + text_font.getsize(now_str+chn)[1] + 15])
        now_str = ''
        split_num = len(now_height) - 1

    try:
        post['text'] = get_content_text(post['text'])
    except Exception as e:
        print('text Error', e)

    for chn in post.get('text'):
        # 每次新增一个文字 如果超过线宽 就把之前的文字存起来
        if chn == '\n' or text_font.getsize(now_str+chn)[0] > text_width*w:
            now_height[-1].append(now_str)
            if not chn == '\n':
                now_str = chn
            else:
                now_str = ''
            now_height.append([now_height[-1][0] + text_font.getsize(now_str+chn)[1] + 10])
        else:
            now_str += chn
    now_height[-1].append(now_str)
    now_height.append([now_height[-1][0] + text_font.getsize(now_str+chn)[1]])

    # 加载博文中的图片
    picHeight = 0  # 记录这些图片总和高度
    pics = []
    for pic in post.get('picUrls'):
        try:
            response = requests.get('https://wx4.sinaimg.cn/large/'+pic+'.jpg', headers=headers)  # 请求图片
            bg = Image.open(BytesIO(response.content))  # 读取图片
            bg = bg.resize((int(text_width*w), int(bg.height*text_width*w/bg.width)), Image.ANTIALIAS)  # 调整大小
            picHeight += bg.height + 3
            pics.append(bg)
        except Exception as e:
            print('[ERROR]: 图片加载错误', 'https://wx4.sinaimg.cn/large/'+pic+'.jpg', e)

    # 测量 发送时间以及设备 文本的高宽
    url_font = ImageFont.truetype(font_type, int(75*text_width*w/tw))
    dt = datetime.strptime(post['time'], '%a %b %d %H:%M:%S %z %Y')
    dt = dt.strftime('%Y-%m-%d %H:%M')[2:]
    s = dt + ' 来自' + post['source']
    tw, th = url_font.getsize(s)

    # 新建底片用以书写内容
    image = Image.new("RGB", (w, int(now_height[-1][0]+0.55*w+picHeight+th+20)), (255, 255, 255))
    draw = ImageDraw.Draw(image)

    # 背景发卡
    faka = Image.open(os.path.join('.', 'gothon', 'card.png'))
    a = Image.new('L', faka.size, 80)  # 透明度 80/255
    i0 = 0
    j0 = 0
    for i in range(i0, ceil(image.height/faka.height)+1, 3):
        for j in range(j0-i, ceil(image.width/faka.width)+1, 4):
            image.paste(faka, (j*faka.width, i*faka.height), mask=a)  # 把发卡贴在背景上

    # 左上角黑框
    draw.rectangle([(0.05*w, 0.28*w), (0.088*w, 0.318*w)], fill='black')
    draw.rectangle([(0.0595*w, 0.2895*w), (0.088*w, 0.318*w)], fill='white')

    # 逐行书写博文
    if post.get('repo'):
        draw.rectangle([
            (int((1-text_width)/2*w-10), int(0.315*w+now_height[split_num][0]-10)),
            (int((1+text_width)/2*w+10), int(0.315*w+now_height[-1][0]))
        ], fill='#f9f9f9')
        draw.line([
            (int((1-text_width)/2*w-10), int(0.315*w+now_height[split_num][0]-10)),
            (int((1+text_width)/2*w+10), int(0.315*w+now_height[split_num][0]-10))
        ], fill='black', width=3)
    for t in now_height[:-1]:
        draw.text((int((1-text_width)/2*w), int(0.315*w+t[0])), t[1], '#343434', text_font)

    # 图片
    now_pic_height = 23
    for pic in pics:
        image.paste(pic, (int((1-text_width)/2*w), int(0.315*w+now_height[-1][0]+now_pic_height)))  # 粘贴至背景
        now_pic_height += pic.height + 3

    # 写发送时间
    draw.text((0.95*w-tw, 0.35*w+now_height[-1][0]+now_pic_height+10), s, '#bebebe', url_font)

    # 二维码区域背景
    draw.rectangle([(0, int(now_height[-1][0]+0.37*w+now_pic_height+th)), (w, int(now_height[-1][0]+0.55*w+now_pic_height+th))], fill='#f9f9f9')

    # 右下角黑框
    draw.rectangle([(0.912*w, 0.312*w+now_height[-1][0]+now_pic_height), (0.95*w, 0.35*w+now_height[-1][0]+now_pic_height)], fill='black')
    draw.rectangle([(0.912*w, 0.312*w+now_height[-1][0]+now_pic_height), (0.9405*w, 0.3405*w+now_height[-1][0]+now_pic_height)], fill='white')

    # 二维码
    url = 'https://m.weibo.com/' + str(userInfo['id']) + '/' + post['mid']
    draw.text((int(0.05*w), int(now_height[-1][0]+0.42*w+now_pic_height+th)), '扫描二维码查看这条动态', '#666666', text_font)
    draw.text((int(0.05*w), int(now_height[-1][0]+0.47*w+now_pic_height+th)), url, '#bebebe', url_font)
    qrimg = qrcode.make(url).resize((int(0.16*w), int(0.16*w)), Image.ANTIALIAS)
    image.paste(qrimg, (int(0.83*w), int(now_height[-1][0]+0.38*w+now_pic_height+th)))

    # 头像
    response = requests.get(userInfo['face'])  # 请求图片
    face = Image.open(BytesIO(response.content))  # 读取图片
    a = Image.new('L', face.size, 0)  # 创建一个黑色背景的画布
    ImageDraw.Draw(a).ellipse((0, 0, a.width, a.height), fill=255)  # 画白色圆形
    face = face.resize((int(0.15*w), int(0.15*w)), Image.ANTIALIAS)
    a = a.resize((int(0.15*w), int(0.15*w)), Image.ANTIALIAS)  # 用于把方形头像变圆
    image.paste(face, (int(0.05*w), int(0.05*w)), mask=a)  # 粘贴至背景

    # 名字
    nw, nh = name_font.getsize(userInfo['name'])
    name_font = ImageFont.truetype(font_blod_type, int(7.5*w/nh))
    draw.text((int(0.21*w), int(0.05*w)), userInfo['name'], '#343434', name_font)

    # 关注与粉丝数
    follow_font = ImageFont.truetype(font_type, int(3.75*w/nh))
    draw.text((int(0.21*w), int(0.14375*w)), '粉丝 {follower} | 关注 {follow}'.format_map(userInfo), '#bebebe', follow_font)

    # 博主个性签名
    desc_font = ImageFont.truetype(font_song, int(12*w/nh))
    draw.text((int(0.015*w), int(0.21*w)), '“', '#787878', desc_font)
    dw, dh = desc_font.getsize('“')
    draw.text((int(0.015*w+dw), int(0.225*w)), userInfo['desc'], '#666666', follow_font)
    return image.save

import aiohttp
import asyncio


a = 10
def SayHello(xixi):
    return xixi + "haha"


async def get(s):
    print(s)
    async with aiohttp.ClientSession() as session:
        resp = await session.get('https://www.baidu.com')
    return resp.status


loop = asyncio.new_event_loop()
res = loop.run_until_complete(get())
loop.run_until_complete(asyncio.sleep(0.3))
print(res)
# By kririae
from PIL import Image
import imagehash
import os, shutil
import math

PATH = ""
FFMPEG_PATH = ""
SIM = 5

class DisjointSet:
  def __init__(self, n: 'int') -> 'void':
    if type(n) != int:
      raise Exception('Input type error!')
    self._fa = list(range(0, n + 1))
  def find(self, x: 'int') -> 'int':
    return x if x == self._fa[x] else self.find(self._fa[x])
  def union(self, x: 'int', y: 'int') -> 'void':
    self._fa[self.find(x)] = self.find(y)

class Img:
  def __init__(self, name):
    img = Image.open(PATH + "\\" + name)
    self.name = name
    self.hash = imagehash.dhash(img)
    self.size = img.size[0] * img.size[1]
    if self.name.endswith('.png'):
      self.size += 1 # correction, Guaranteed that under the same quality, png first.

def extractIntegers(s: 'str') -> 'list':
  if type(s) != str:
    raise Exception('Input type error!')
  s += '#'
  ret, val = [], 0
  for index, value in enumerate(s):
    if value.isdigit():
      val = val * 10 + int(value)
    elif s[index - 1].isdigit() if index >= 1 else False:
      ret.append(val)
      val = 0
  return ret

def formatName(img: 'Img') -> 'str':
  s = img.name
  # a = lambda x: ''.join(x.split('.')[:-1]) + '.jpg'
  end = lambda x: x.split('.')[-1]
  if 'Konachan.com' in s:
    return f'Konachan.com-{extractIntegers(s)[0]}.{end(s)}'
  if '_p0' in s:
    return s
  if 'yande.re' in s:
    return f'yande.re-{str(extractIntegers(s)[0])}.{end(s)}'
  return f'{str(img.hash)}.{end(s)}'

def isImage(x: 'string') -> 'bool':
  if type(x) != str:
    raise Exception('Input type error!')
  return x.endswith('.jpg') or x.endswith('.png')

def transformImage(_from: 'str', _to: 'str') -> 'void':
  # reduce its quality...
  if type(_from) != str or type(_to) != str:
    raise Exception('Input type error!')
  img = Image.open(_from)
  img.save(to)

if __name__ == '__main__':
  _dir = list(filter(isImage, os.listdir(PATH)))
  imgs = []
  for i in range(0, len(_dir)):
    imgs.append(Img(_dir[i]))
  unq = DisjointSet(len(_dir))
  for i in range(0, len(_dir)):
    for j in range(i + 1, len(_dir)):
      if abs(imgs[i].hash - imgs[j].hash) <= SIM:
        unq.union(i + 1, j + 1)

  ref = {}
  for i in range(1, len(unq._fa)):
    rt = unq.find(i)
    if ref.get(rt) != None:
      ref[rt].append(i - 1)
    else:
      ref[rt] = [i - 1]
  for i, j in ref.items():
    ref[i] = sorted(j, key=lambda x: imgs[x].size, reverse=True)

  if not os.path.exists(PATH + '\\' + 'filted'):
    os.mkdir(PATH + '\\' + 'filted')
  for i, j in ref.items():
    hi_quality_img = imgs[j[0]]
    if len(j) > 1:
      print(f'Source: {hi_quality_img.name}')
      for k in range(1, len(j)):
        print(f'---- {imgs[j[k]].name}')
    fr = PATH + '\\' + hi_quality_img.name
    _name = formatName(hi_quality_img)
    to = f'{PATH}\\filted\\{_name}'
    shutil.copyfile(fr, to)
    # transformImage(fr, to)

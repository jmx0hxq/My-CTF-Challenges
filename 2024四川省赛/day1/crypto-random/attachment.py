import random
from Crypto.Util.number import *

with open('flag') as (f):
    flag = list(f.read())
p=getPrime(23)
g=2
a=random.randint(g,p-1)
b=random.randint(g,p-1)
k=pow(g,a*b,p)
print(p)

random.seed(k)
for i in range(0,len(flag),2):
    flag[i],flag[i+1]=flag[i+1],flag[i]
else:
    random.shuffle(flag)
print(flag)


# 5037523
# ['1', '0', 'a', '}', 'b', '0', '9', 'c', 'b', '9', 'a', 'g', '-', '-', 'c', '0', '1', 'c', '9', '4', 'e', '3', '{', 'f', 'f', '7', '1', 'c', 'l', '0', '7', 'e', '6', '9', '8', 'b', '7', '-', 'a', '-', 'f', '3']
import os
import re

# 从环境变量获取行号并增加3
line = int(os.environ.get('GOLINE')) + 3

# 使用utf-8编码读取README.md文件
with open('README.md', 'r', encoding='utf-8') as f:
    readme = f.read()

# 使用正则表达式查找包含特定链接的行
link = re.search(r'https://github.com/2404589803/deepspace/blob/main/persistence\.go#L(\d+)', readme)

# 如果找到匹配的链接，则进行替换
if link:
    replaced = readme[:link.start()] + f'https://github.com/2404589803/deepspace/blob/main/persistence.go#L{line}' + readme[link.end():]

    # 使用utf-8编码写入修改后的内容到README.md文件
    with open('README.md', 'w', encoding='utf-8') as f:
        f.write(replaced)
else:
    print("No matching link found in README.md.")

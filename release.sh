#!/bin/bash

# 获取当前最新的标签
latest_tag=$(git describe --tags --abbrev=0)

# 提取标签中的最后一组数字
tag_number=$(echo "$latest_tag" | grep -oE '[0-9]+$')

# 提取标签的前缀（包含可能的数字）
prefix=$(echo "$latest_tag" | sed -E "s/$tag_number\$//")

# 将数字部分加1
new_tag_number=$((tag_number + 1))

# 构造新的标签
new_tag="${prefix}${new_tag_number}"

# 提示用户确认
read -p "新标签：${new_tag}，是否确认创建并推送新标签？ (y/n): " confirm

if [ "$confirm" == "y" ]; then
    # 打印新标签并加到仓库
    echo "正在创建并推送新标签：$new_tag"
    git tag "$new_tag"
    git push origin "$new_tag"
    echo "新标签已成功创建并推送。"
else
    echo "操作已取消。"
fi
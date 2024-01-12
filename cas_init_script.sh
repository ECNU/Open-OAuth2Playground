#!/bin/bash
# cas_init_script.sh

mkdir -p /export/data/
chmod 777 /export/data/
sqlite3 /export/data/cas.db <<EOF
CREATE TABLE IF NOT EXISTS user (username TEXT, password TEXT, name TEXT);
DELETE FROM user;
INSERT INTO user (username, password, name) VALUES ('cas', '123456', '测试用户');
EOF

echo "cas.db created successfully!"

# 读取环境变量，如果未设置，则使用默认值
CAS_SERVER_NAME=${CAS_SERVER_NAME:-"http://localhost:8444"}
SERVER_PORT=${SERVER_PORT:-"8444"}

# 配置文件路径
CAS_PROPERTIES_FILE="/etc/cas/config/cas.properties"

# 检查并替换或添加 server.port
if grep -q "server.port" "$CAS_PROPERTIES_FILE"; then
    sed -i "s#server.port=.*#server.port=${SERVER_PORT}#" "$CAS_PROPERTIES_FILE"
else
    echo "server.port=${SERVER_PORT}" >> "$CAS_PROPERTIES_FILE"
fi

# 检查并替换或添加 cas.server.name
if grep -q "cas.server.name" "$CAS_PROPERTIES_FILE"; then
    sed -i "s#cas.server.name=.*#cas.server.name=${CAS_SERVER_NAME}#" "$CAS_PROPERTIES_FILE"
else
    echo "cas.server.name=${CAS_SERVER_NAME}" >> "$CAS_PROPERTIES_FILE"
fi

echo "read configuration successfully!"
#!/bin/sh
set -e

echo "[Finlog] 🔧 Starting custom Netdata initialization..."

CONF="/etc/netdata/go.d/mysql.conf"
echo "[Finlog] 📄 Generating $CONF..."

cat > "$CONF" <<EOF
jobs:
  - name: finlog-db
    dsn: ${NETDATA_MYSQL_USER}:${NETDATA_MYSQL_PASSWORD}@tcp(finlog-db:3306)/
EOF

chmod 644 "$CONF"

echo "[Finlog] 🔧 Enabling cgroups plugin..."
if [ -f /etc/netdata/netdata.conf ]; then
    if ! grep -q "\[plugin:cgroups\]" /etc/netdata/netdata.conf; then
        echo -e "\n[plugin:cgroups]\n  update every = 1\n  command options = -r" >> /etc/netdata/netdata.conf
    fi
fi

echo "[Finlog] 🔧 Enabling Docker monitoring..."
if [ -f /etc/netdata/go.d.conf ]; then
    sed -i 's/# docker: no/docker: yes/' /etc/netdata/go.d.conf || true
fi

echo "[Finlog] ⏳ Waiting for MySQL (finlog-db:3306)..."
MAX_RETRIES=30
RETRY_COUNT=0

while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
    if timeout 5 bash -c 'cat < /dev/null > /dev/tcp/finlog-db/3306' 2>/dev/null; then
        echo "[Finlog] ✅ MySQL is ready!"
        
        echo "[Finlog] 🔍 Testing MySQL connection..."
        if mysql -h finlog-db -u "${NETDATA_MYSQL_USER}" -p"${NETDATA_MYSQL_PASSWORD}" -e "SELECT 1" 2>/dev/null; then
            echo "[Finlog] ✅ MySQL connection successful!"
        else
            echo "[Finlog] ⚠️  MySQL credentials might be incorrect"
        fi
        break
    fi
    RETRY_COUNT=$((RETRY_COUNT + 1))
    echo "[Finlog] ⏳ Waiting for MySQL... ($RETRY_COUNT/$MAX_RETRIES)"
    sleep 2
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "[Finlog] ⚠️  MySQL not ready after 60s, continuing anyway..."
fi

echo "[Finlog] 🚀 Starting Netdata..."
exec /usr/sbin/run.sh "$@"
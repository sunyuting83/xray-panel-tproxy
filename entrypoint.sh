#!/bin/sh
set -e

check_and_delete_rule() {
    if ip rule show | grep -q "fwmark 0x1 lookup 100"; then
        ip rule delete fwmark 1 table 100
        ip route delete local default dev lo table 100
    fi
    if ip -6 rule show | grep -q "fwmark 0x1 lookup 106"; then
        ip -6 rule delete fwmark 1 table 106
        ip -6 route delete local ::/0 dev lo table 106
    fi
}

reset_iptables(){
    echo "Resetting iptables rules..."
    check_and_delete_rule

    iptables -t mangle -D PREROUTING -j XRAY 2>/dev/null || true
    iptables -t mangle -D PREROUTING -p tcp -m socket -j DIVERT 2>/dev/null || true
    iptables -t mangle -F XRAY 2>/dev/null || true
    iptables -t mangle -X XRAY 2>/dev/null || true
    iptables -t mangle -F DIVERT 2>/dev/null || true
    iptables -t mangle -X DIVERT 2>/dev/null || true

    ip6tables-nft -t mangle -D PREROUTING -j XRAY6_MASK 2>/dev/null || true
    ip6tables-nft -t mangle -D PREROUTING -p tcp -m socket -j DIVERT 2>/dev/null || true
    ip6tables-nft -t mangle -F XRAY6_MASK 2>/dev/null || true
    ip6tables-nft -t mangle -X XRAY6_MASK 2>/dev/null || true
    ip6tables-nft -t mangle -F DIVERT 2>/dev/null || true
    ip6tables-nft -t mangle -X DIVERT 2>/dev/null || true
}

set_xray_iptables(){
    echo "Setting up Xray iptables rules..."

    ip rule add fwmark 1 table 100 2>/dev/null || true
    ip route add local default dev lo table 100 2>/dev/null || true
    ip -6 rule add fwmark 1 table 106 2>/dev/null || true
    ip -6 route add local ::/0 dev lo table 106 2>/dev/null || true

    iptables -t mangle -N XRAY
    iptables -t mangle -A XRAY -d 127.0.0.0/8 -j RETURN
    iptables -t mangle -A XRAY -d 10.0.0.0/8 -j RETURN
    iptables -t mangle -A XRAY -d 172.16.0.0/12 -j RETURN
    iptables -t mangle -A XRAY -d 192.168.0.0/16 -j RETURN
    iptables -t mangle -A XRAY -d 192.168.0.0/16 -p tcp ! --dport 53 -j RETURN
    iptables -t mangle -A XRAY -d 192.168.0.0/16 -p udp ! --dport 53 -j RETURN
    iptables -t mangle -A XRAY -j RETURN -m mark --mark 0xff
    iptables -t mangle -A XRAY -p tcp -j TPROXY --on-port 7892 --tproxy-mark 1
    iptables -t mangle -A XRAY -p udp -j TPROXY --on-port 7892 --tproxy-mark 1
    iptables -t mangle -A PREROUTING -j XRAY

    ip6tables-nft -t mangle -N XRAY6_MASK
    ip6tables-nft -t mangle -A XRAY6_MASK -d fe80::/10 -j RETURN
    ip6tables-nft -t mangle -A XRAY6_MASK -d fd00::/8 -j RETURN
    ip6tables-nft -t mangle -A XRAY6_MASK -d fd00::/8 -p tcp ! --dport 53 -j RETURN
    ip6tables-nft -t mangle -A XRAY6_MASK -d fd00::/8 -p udp ! --dport 53 -j RETURN
    ip6tables-nft -t mangle -A XRAY6_MASK -j RETURN -m mark --mark 0xff
    ip6tables-nft -t mangle -A XRAY6_MASK -p udp -j TPROXY --on-port 7892 --tproxy-mark 1
    ip6tables-nft -t mangle -A XRAY6_MASK -p tcp -j TPROXY --on-port 7892 --tproxy-mark 1
    ip6tables-nft -t mangle -A PREROUTING -j XRAY6_MASK

    iptables -t mangle -N DIVERT
    iptables -t mangle -A DIVERT -j MARK --set-mark 1
    iptables -t mangle -A DIVERT -j ACCEPT
    iptables -t mangle -I PREROUTING -p tcp -m socket -j DIVERT

    ip6tables-nft -t mangle -N DIVERT
    ip6tables-nft -t mangle -A DIVERT -j MARK --set-mark 1
    ip6tables-nft -t mangle -A DIVERT -j ACCEPT
    ip6tables-nft -t mangle -I PREROUTING -p tcp -m socket -j DIVERT

    iptables -t mangle -I XRAY 1 -i docker0 -j RETURN
    ip6tables-nft -t mangle -I XRAY6_MASK 1 -i docker0 -j RETURN
}

reset_iptables
set_xray_iptables

echo "Iptables applied successfully."

/xpanel/server
exec "$@"
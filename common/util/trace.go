package util

import (
	"context"
	"encoding/binary"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

func GenerateSpanID(addr string) string {
	strAddr := strings.Split(addr, ":")
	ip := strAddr[0]
	ipLong, _ := Ip2Long(ip)
	times := uint64(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	spanId := ((times ^ uint64(ipLong)) << 32) | uint64(rand.Int31())
	return strconv.FormatUint(spanId, 16)
}

func Ip2Long(ip string) (uint32, error) {
	ipAddr, err := net.ResolveIPAddr("ip", ip)
	if err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint32(ipAddr.IP.To4()), nil
}

func GetTraceInfoFromCtx(ctx context.Context) (traceId, spanId, pSpanId string) {
	if ctx.Value("traceid") != nil {
		traceId = ctx.Value("traceid").(string)
	}
	if ctx.Value("spanid") != nil {
		spanId = ctx.Value("spanid").(string)
	}
	if ctx.Value("psapnid") != nil {
		pSpanId = ctx.Value("pspanid").(string)
	}

	return
}

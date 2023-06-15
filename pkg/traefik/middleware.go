package traefik

import (
	"strconv"
	"strings"
)

type HttpMiddleware struct {
	name string
	kvs  map[string]string
}

func (m HttpMiddleware) Name() string {
	return m.name
}
func (m HttpMiddleware) GetKvs() map[string]string {
	return m.kvs
}

type TcpMiddleware struct {
	name string
	kvs  map[string]string
}

func (m TcpMiddleware) Name() string {
	return m.name
}
func (m TcpMiddleware) GetKvs() map[string]string {
	return m.kvs
}

// NewHttpAddPrefix 添加前缀 traefik/http/middlewares/Middleware00/addPrefix/prefix /foo
func NewHttpAddPrefix(name string, prefix string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	if prefix != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "addPrefix/prefix")] = "/" + strings.TrimPrefix(prefix, "/")
	}
	return
}

/*
traefik/http/middlewares/Middleware21/stripPrefix/forceSlash	true
traefik/http/middlewares/Middleware21/stripPrefix/prefixes/0	foobar
traefik/http/middlewares/Middleware21/stripPrefix/prefixes/1	foobar
*/

// NewHttpStripPrefix x
func NewHttpStripPrefix(name string, prefixes []string, forceSlash bool) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	for i, p := range prefixes {
		if p != "" {
			m.kvs[EtcdKey(httpMidPrefix, "stripPrefix/prefixes", strconv.Itoa(i))] = "/" + strings.TrimPrefix(p, "/")
		}
	}
	m.kvs[EtcdKey(httpMidPrefix, name, "stripPrefix/prefixes/forceSlash")] = BoolVal(forceSlash)
	return
}

/*
traefik/http/middlewares/Middleware01/basicAuth/headerField	foobar
traefik/http/middlewares/Middleware01/basicAuth/realm	foobar
traefik/http/middlewares/Middleware01/basicAuth/removeHeader	true
traefik/http/middlewares/Middleware01/basicAuth/users/0	foobar
traefik/http/middlewares/Middleware01/basicAuth/users/1	foobar
traefik/http/middlewares/Middleware01/basicAuth/usersFile	foobar
*/

// NewHttpBasicAuth auth basic认证 users= {user:password} realm=MyRealm headerField=X-WebAuth-User
func NewHttpBasicAuth(name string, users map[string]string, headerField, realm string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	i := 0
	for acc, pwd := range users {
		m.kvs[EtcdKey(httpMidPrefix, name, "basicAuth/users", strconv.Itoa(i))] = acc + ":" + pwd
		i++
	}
	if headerField != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "basicAuth/headerField")] = headerField
	}
	if realm != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "basicAuth/realm")] = realm
	}
	m.kvs[EtcdKey(httpMidPrefix, name, "basicAuth/removeHeader")] = "true"
	return
}

/* file content：
test:$apr1$H6uskkkW$IgXLP6ewTrSuBkTrqE8wj/
test2:$apr1$d9hr9HBB$4HxwgUir3HP4EsggP/QNo0
*/

// NewHttpBasicAuthWithFile 文件提供user的auth basic 认证 users in file
func NewHttpBasicAuthWithFile(name string, usersFile string, headerField, realm string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	if headerField != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "basicAuth/headerField")] = headerField
	}
	if realm != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "basicAuth/realm")] = realm
	}
	m.kvs[EtcdKey(httpMidPrefix, name, "basicAuth/removeHeader")] = "true"
	m.kvs[EtcdKey(httpMidPrefix, name, "basicAuth/usersFile")] = usersFile
	return
}

/*
traefik/http/middlewares/Middleware02/buffering/maxRequestBodyBytes	42
traefik/http/middlewares/Middleware02/buffering/maxResponseBodyBytes	42
traefik/http/middlewares/Middleware02/buffering/memRequestBodyBytes	42
traefik/http/middlewares/Middleware02/buffering/memResponseBodyBytes	42
traefik/http/middlewares/Middleware02/buffering/retryExpression	foobar
*/

// NewHttpBasicBuffering 限制请求的大小 retryExp重试机制=IsNetworkError() && Attempts() < 2 &&ResponseCode()=xx
func NewHttpBasicBuffering(name string, maxReqBodyBytes, maxRespBodyBytes, memReqBodyBytes, memResBodyBytes int, retryExp string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	if maxReqBodyBytes > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "buffering/maxRequestBodyBytes")] = strconv.Itoa(maxReqBodyBytes)
	}
	if maxRespBodyBytes > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "buffering/maxResponseBodyBytes")] = strconv.Itoa(maxRespBodyBytes)
	}
	if memReqBodyBytes > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "buffering/memRequestBodyBytes")] = strconv.Itoa(memReqBodyBytes)
	}
	if memResBodyBytes > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "buffering/memResponseBodyBytes")] = strconv.Itoa(memResBodyBytes)
	}
	if retryExp != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "buffering/retryExpression")] = retryExp
	}
	return
}

/*
traefik/http/middlewares/Middleware03/chain/middlewares/0	foobar
traefik/http/middlewares/Middleware03/chain/middlewares/1	foobar
*/

// NewHttpChain 中间件组
func NewHttpChain(name string, middlewares []string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	for i, mid := range middlewares {
		if mid != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "chain/middlewares", strconv.Itoa(i))] = mid
		}
	}
	return
}

/*
traefik/http/middlewares/Middleware04/circuitBreaker/checkPeriod	42s 检查间隔期
traefik/http/middlewares/Middleware04/circuitBreaker/expression	foobar open的触发条件： 网络错误率 NetworkErrorRatio() > 0.30， 错误码：ResponseCodeRatio(500, 600, 0, 600) > 0.25， 变慢：LatencyAtQuantileMS(50.0) > 100
traefik/http/middlewares/Middleware04/circuitBreaker/fallbackDuration	42s 回退机制持续
traefik/http/middlewares/Middleware04/circuitBreaker/recoveryDuration	42s 恢复机制持续
*/

// NewHttpCircuitBreaker 断路器 避免请求堆叠到不正常的服务器上， 服务不正常时打开，将请求转到回退机制
func NewHttpCircuitBreaker(name string, expression string, checkPeriod int, fallbackDuration int, recoveryDuration int) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	if expression != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "circuitBreaker/expression")] = expression
	}

	if checkPeriod > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "circuitBreaker/checkPeriod")] = strconv.Itoa(checkPeriod) + "ms"
	}

	if fallbackDuration > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "circuitBreaker/fallbackDuration")] = strconv.Itoa(fallbackDuration) + "s"
	}

	if recoveryDuration > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "circuitBreaker/recoveryDuration")] = strconv.Itoa(recoveryDuration) + "s"
	}
	return
}

/*

traefik/http/middlewares/Middleware05/compress/excludedContentTypes/0	foobar
traefik/http/middlewares/Middleware05/compress/excludedContentTypes/1	foobar
traefik/http/middlewares/Middleware05/compress/minResponseBodyBytes	42
*/

// NewHttpCompress 压缩
func NewHttpCompress(name string, minResponseBodyBytes int, excludedContentTypes []string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	if minResponseBodyBytes > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "compress/minResponseBodyBytes")] = strconv.Itoa(minResponseBodyBytes)
	}
	for i, ct := range excludedContentTypes {
		m.kvs[EtcdKey(httpMidPrefix, name, "compress/excludedContentTypes", strconv.Itoa(i))] = ct
	}
	return
}

// NewHttpDetectContentType content-type后端未设置时是否从内容自动设置 traefik/http/middlewares/Middleware06/contentType/autoDetect	true
func NewHttpDetectContentType(name string, enable bool) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	m.kvs[EtcdKey(httpMidPrefix, name, "compress/contentType/autoDetect")] = BoolVal(enable)
	return
}

/*
traefik/http/middlewares/Middleware08/errors/query	foobar
traefik/http/middlewares/Middleware08/errors/service	foobar
traefik/http/middlewares/Middleware08/errors/status/0	foobar
traefik/http/middlewares/Middleware08/errors/status/1	foobar
*/

// NewHttpErrPage 自定义错误页面 statusCode="500-599" "401,402,500-503", page="{status}.html"
func NewHttpErrPage(name string, statusCodes []string, serviceName string, page string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	m.kvs[EtcdKey(httpMidPrefix, name, "errors/service")] = serviceName
	m.kvs[EtcdKey(httpMidPrefix, name, "errors/query")] = page
	for i, s := range statusCodes {
		m.kvs[EtcdKey(httpMidPrefix, name, "errors/status", strconv.Itoa(i))] = s
	}
	return
}

type HeaderOptions struct {
	RequestHeaders   map[string]string // 添加或删除的请求头 空则是删除
	ResponseHeaders  map[string]string // 添加或删除的响应头 为空则是删除
	frameDeny        bool              // add X-Frame-Options 头
	browserXssFilter bool              // X-XSS-Protection
	AddVertHeader    bool
	Cors             CORSOptions
}
type CORSOptions struct {
	AllowOrigins      []string
	AllowRegexOrigins []string
	AllowHeaders      []string
	ExposeHeaders     []string
	AllowMethods      []string
	MaxAge            int
	AllowCredentials  bool
}

/*
traefik/http/middlewares/Middleware10/headers/frameDeny	true
traefik/http/middlewares/Middleware10/headers/browserXssFilter	true
traefik/http/middlewares/Middleware10/headers/customRequestHeaders/name0	foobar
traefik/http/middlewares/Middleware10/headers/customRequestHeaders/name1	foobar
traefik/http/middlewares/Middleware10/headers/customResponseHeaders/name0	foobar
traefik/http/middlewares/Middleware10/headers/customResponseHeaders/name1	foobar

traefik/http/middlewares/Middleware10/headers/accessControlAllowCredentials	true
traefik/http/middlewares/Middleware10/headers/accessControlAllowHeaders/0	foobar
traefik/http/middlewares/Middleware10/headers/accessControlAllowHeaders/1	foobar
traefik/http/middlewares/Middleware10/headers/accessControlAllowMethods/0	foobar
traefik/http/middlewares/Middleware10/headers/accessControlAllowMethods/1	foobar
traefik/http/middlewares/Middleware10/headers/accessControlAllowOriginList/0	foobar
traefik/http/middlewares/Middleware10/headers/accessControlAllowOriginList/1	foobar
traefik/http/middlewares/Middleware10/headers/accessControlAllowOriginListRegex/0	foobar
traefik/http/middlewares/Middleware10/headers/accessControlAllowOriginListRegex/1	foobar
traefik/http/middlewares/Middleware10/headers/accessControlExposeHeaders/0	foobar
traefik/http/middlewares/Middleware10/headers/accessControlExposeHeaders/1	foobar
traefik/http/middlewares/Middleware10/headers/accessControlMaxAge	42
traefik/http/middlewares/Middleware10/headers/addVaryHeader	true
*/

// NewHttpHeader 请求和响应的头字段管理 add or remove  为空则时删除
func NewHttpHeader(name string, o HeaderOptions) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	for k, v := range o.RequestHeaders {
		if k != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "headers/customRequestHeaders", k)] = v
		}
	}
	for k, v := range o.ResponseHeaders {
		if k != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "headers/customResponseHeaders", k)] = v
		}
	}
	m.kvs[EtcdKey(httpMidPrefix, name, "headers/customResponseHeaders/frameDeny")] = BoolVal(o.frameDeny)
	m.kvs[EtcdKey(httpMidPrefix, name, "headers/customResponseHeaders/browserXssFilter")] = BoolVal(o.browserXssFilter)
	m.kvs[EtcdKey(httpMidPrefix, name, "headers/accessControlAllowCredentials")] = BoolVal(o.Cors.AllowCredentials)
	m.kvs[EtcdKey(httpMidPrefix, name, "headers/addVaryHeader")] = BoolVal(o.AddVertHeader)

	for k, v := range o.Cors.AllowMethods {
		m.kvs[EtcdKey(httpMidPrefix, name, "headers/accessControlAllowMethods", strconv.Itoa(k))] = v
	}
	for k, v := range o.Cors.AllowHeaders {
		m.kvs[EtcdKey(httpMidPrefix, name, "headers/accessControlAllowHeaders", strconv.Itoa(k))] = v
	}
	for k, v := range o.Cors.AllowOrigins {
		m.kvs[EtcdKey(httpMidPrefix, name, "headers/accessControlAllowOriginList", strconv.Itoa(k))] = v
	}
	for k, v := range o.Cors.AllowRegexOrigins {
		m.kvs[EtcdKey(httpMidPrefix, name, "headers/accessControlAllowOriginListRegex", strconv.Itoa(k))] = v
	}
	for k, v := range o.Cors.ExposeHeaders {
		m.kvs[EtcdKey(httpMidPrefix, name, "headers/accessControlExposeHeaders", strconv.Itoa(k))] = v
	}
	if o.Cors.MaxAge > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "headers/accessControlMaxAge")] = strconv.Itoa(o.Cors.MaxAge)
	}
	return
}

/*
traefik/http/middlewares/Middleware11/ipWhiteList/ipStrategy/depth	42
traefik/http/middlewares/Middleware11/ipWhiteList/ipStrategy/excludedIPs/0	foobar
traefik/http/middlewares/Middleware11/ipWhiteList/ipStrategy/excludedIPs/1	foobar
traefik/http/middlewares/Middleware11/ipWhiteList/sourceRange/0	foobar
traefik/http/middlewares/Middleware11/ipWhiteList/sourceRange/1	foobar
*/

// NewHttpWhitelist ip白名单  depth, X-Forwarded-For
func NewHttpWhitelist(name string, excludeIps []string, depth int, ips []string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	if depth > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "ipWhiteList/ipStrategy/depth")] = strconv.Itoa(depth)
	}
	for i, ip := range excludeIps {
		if ip != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "ipWhiteList/ipStrategy/excludedIPs", strconv.Itoa(i))] = ip
		}
	}
	for i, ip := range ips {
		if ip != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "ipWhiteList/sourceRange", strconv.Itoa(i))] = ip
		}
	}
	return
}

/*
traefik/http/middlewares/Middleware12/inFlightReq/amount	42
traefik/http/middlewares/Middleware12/inFlightReq/sourceCriterion/ipStrategy/depth	42
traefik/http/middlewares/Middleware12/inFlightReq/sourceCriterion/ipStrategy/excludedIPs/0	foobar
traefik/http/middlewares/Middleware12/inFlightReq/sourceCriterion/ipStrategy/excludedIPs/1	foobar
traefik/http/middlewares/Middleware12/inFlightReq/sourceCriterion/requestHeaderName	foobar
traefik/http/middlewares/Middleware12/inFlightReq/sourceCriterion/requestHost	true
*/

type SourceCriterion struct {
	Depth      int
	ExcludeIps []string
	Header     string
	Host       bool
}

// NewHttpInFlightReq 同时进行的连接限制, amount是数量， depth ips是限制策略
func NewHttpInFlightReq(name string, amount int, criterion SourceCriterion) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	if amount > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "inFlightReq/amount")] = strconv.Itoa(amount)
	}
	if criterion.Depth > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "inFlightReq/sourceCriterion/ipStrategy/depth")] = strconv.Itoa(criterion.Depth)
	}
	if criterion.Header != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "inFlightReq/sourceCriterion/requestHeaderName")] = criterion.Header
	}
	if criterion.Host {
		m.kvs[EtcdKey(httpMidPrefix, name, "inFlightReq/sourceCriterion/requestHost")] = "true"
	}
	for i, ip := range criterion.ExcludeIps {
		if ip != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "inFlightReq/sourceCriterion/ipStrategy/excludedIPs", strconv.Itoa(i))] = ip
		}
	}
	return
}

/*
traefik/http/middlewares/Middleware15/rateLimit/average	42
traefik/http/middlewares/Middleware15/rateLimit/burst	42
traefik/http/middlewares/Middleware15/rateLimit/period	42s
traefik/http/middlewares/Middleware15/rateLimit/sourceCriterion/ipStrategy/depth	42
traefik/http/middlewares/Middleware15/rateLimit/sourceCriterion/ipStrategy/excludedIPs/0	foobar
traefik/http/middlewares/Middleware15/rateLimit/sourceCriterion/ipStrategy/excludedIPs/1	foobar
traefik/http/middlewares/Middleware15/rateLimit/sourceCriterion/requestHeaderName	foobar
traefik/http/middlewares/Middleware15/rateLimit/sourceCriterion/requestHost	true
*/

type LimitConf struct {
	Limit  int
	Period int
	Max    int
}

// NewHttpRateLimit 限流 limit reqs / period seconds
func NewHttpRateLimit(name string, conf LimitConf, criterion SourceCriterion) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	if conf.Limit > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "rateLimit/average")] = strconv.Itoa(conf.Limit)
	}
	if conf.Period > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "rateLimit/period")] = strconv.Itoa(conf.Period) + "s"
	}
	if conf.Max > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "rateLimit/burst")] = strconv.Itoa(conf.Max)
	}

	if criterion.Depth > 0 {
		m.kvs[EtcdKey(httpMidPrefix, name, "rateLimit/sourceCriterion/ipStrategy/depth")] = strconv.Itoa(criterion.Depth)
	}
	if criterion.Header != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "rateLimit/sourceCriterion/requestHeaderName")] = criterion.Header
	}
	if criterion.Host {
		m.kvs[EtcdKey(httpMidPrefix, name, "rateLimit/sourceCriterion/requestHost")] = "true"
	}
	for i, ip := range criterion.ExcludeIps {
		if ip != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "rateLimit/sourceCriterion/ipStrategy/excludedIPs", strconv.Itoa(i))] = ip
		}
	}
	return
}

/*
traefik/http/middlewares/Middleware09/forwardAuth/address	foobar
traefik/http/middlewares/Middleware09/forwardAuth/authRequestHeaders/0	foobar
traefik/http/middlewares/Middleware09/forwardAuth/authRequestHeaders/1	foobar
traefik/http/middlewares/Middleware09/forwardAuth/authResponseHeaders/0	foobar
traefik/http/middlewares/Middleware09/forwardAuth/authResponseHeaders/1	foobar
traefik/http/middlewares/Middleware09/forwardAuth/authResponseHeadersRegex	foobar
traefik/http/middlewares/Middleware09/forwardAuth/tls/ca	foobar
traefik/http/middlewares/Middleware09/forwardAuth/tls/caOptional	true
traefik/http/middlewares/Middleware09/forwardAuth/tls/cert	foobar
traefik/http/middlewares/Middleware09/forwardAuth/tls/insecureSkipVerify	true
traefik/http/middlewares/Middleware09/forwardAuth/tls/key	foobar
traefik/http/middlewares/Middleware09/forwardAuth/trustForwardHeader	true
*/

type AuthConf struct {
	Address              string // https://example.com/auth
	RequestHeaders       []string
	ResponseHeaders      []string
	ResponseHeadersRegex string
	Tls                  *AuthTls
	TrustForwardHeader   bool // 信任转发的头 X-Forwarded-* （Method Proto Host Uri For）
}
type AuthTls struct {
	Ca                 string //"path/to/local.crt"
	CaOptional         bool
	Cert               string // "path/to/foo.cert"
	Key                string // "path/to/foo.key"
	InsecureSkipVerify bool   // 接收任何证书
}

// NewHttpAuth 认证中间件
func NewHttpAuth(name string, conf AuthConf) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/address")] = conf.Address
	for i, k := range conf.RequestHeaders {
		m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/authRequestHeaders", strconv.Itoa(i))] = k
	}
	for i, k := range conf.ResponseHeaders {
		m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/authResponseHeaders", strconv.Itoa(i))] = k
	}
	if conf.ResponseHeadersRegex != "" {
		m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/authResponseHeadersRegex")] = conf.ResponseHeadersRegex
	}
	m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/trustForwardHeader")] = "true"

	if conf.Tls != nil {
		if conf.Tls.Ca != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/tls/ca")] = conf.Tls.Ca
		}
		if conf.Tls.Cert != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/tls/cert")] = conf.Tls.Cert
		}
		if conf.Tls.Key != "" {
			m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/tls/key")] = conf.Tls.Key
		}
		m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/tls/caOptional")] = BoolVal(conf.Tls.CaOptional)
		m.kvs[EtcdKey(httpMidPrefix, name, "forwardAuth/tls/insecureSkipVerify")] = BoolVal(conf.Tls.InsecureSkipVerify)
	}
	return
}

// NewHttpRexReplace 正则替换
func NewHttpRexReplace(name string, regex, replacement string) (m HttpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	m.kvs[EtcdKey(httpMidPrefix, name, "replacePathRegex/regex")] = regex
	m.kvs[EtcdKey(httpMidPrefix, name, "replacePathRegex/replacement")] = replacement
	return
}

// NewTcpWhitelist traefik/tcp/middlewares/TCPMiddleware00/ipWhiteList/sourceRange/0	foobar
func NewTcpWhitelist(name string, ips []string) (m TcpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	for i, v := range ips {
		if v != "" {
			m.kvs[EtcdKey(tcpMidPrefix, name, "ipWhiteList/sourceRange", strconv.Itoa(i))] = v
		}
	}
	return
}

// NewTcpInFlightConn 限制一个端同时进行的连接 traefik/tcp/middlewares/TCPMiddleware01/inFlightConn/amount	42
func NewTcpInFlightConn(name string, amount int) (m TcpMiddleware) {
	m.name = name
	m.kvs = make(map[string]string)
	m.kvs[EtcdKey(tcpMidPrefix, name, "inFlightConn/amount")] = strconv.Itoa(amount)
	return
}

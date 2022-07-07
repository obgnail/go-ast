package decl

type ImportType = string

const (
	dotImport        ImportType = "dot"
	normalImport     ImportType = "normal"
	aliasImport      ImportType = "alias"
	underscoreImport ImportType = "underscore"
)

var (
	StdLibs = []string{
		"archive/tar", "archive/zip", "bufio", "builtin", "bytes",
		"cmd/addr2line", "cmd/api", "cmd/asm", "cmd/buildid", "cmd/cgo",
		"cmd/compile", "cmd/cover", "cmd/dist", "cmd/doc", "cmd/fix",
		"cmd/go", "cmd/gofmt", "cmd/link", "cmd/nm", "cmd/objdump",
		"cmd/pack", "cmd/pprof", "cmd/trace", "cmd/vet", "compress/bzip2",
		"compress/flate", "compress/gzip", "compress/lzw", "compress/zlib", "container/heap",
		"container/list", "container/ring", "context", "crypto", "crypto/aes",
		"crypto/cipher", "crypto/des", "crypto/dsa", "crypto/ecdsa", "crypto/ed25519",
		"crypto/elliptic", "crypto/hmac", "crypto/md5", "crypto/rand", "crypto/rc4",
		"crypto/rsa", "crypto/sha1", "crypto/sha256", "crypto/sha512", "crypto/subtle",
		"crypto/tls", "crypto/x509", "crypto/x509/pkix", "database/sql", "database/sql/driver",
		"debug/dwarf", "debug/elf", "debug/gosym", "debug/macho", "debug/pe",
		"debug/plan9obj", "embed", "encoding", "encoding/ascii85", "encoding/asn1",
		"encoding/base32", "encoding/base64", "encoding/binary", "encoding/csv", "encoding/gob",
		"encoding/hex", "encoding/json", "encoding/pem", "encoding/xml", "errors",
		"expvar", "flag", "fmt", "go/ast", "go/build",
		"go/build/constraint", "go/constant", "go/doc", "go/format", "go/importer",
		"go/parser", "go/printer", "go/scanner", "go/token", "go/types",
		"hash", "hash/adler32", "hash/crc32", "hash/crc64", "hash/fnv",
		"hash/maphash", "html", "html/template", "image", "image/color",
		"image/color/palette", "image/draw", "image/gif", "image/jpeg", "image/png",
		"index/suffixarray", "io", "io/fs", "io/ioutil", "log",
		"log/syslog", "math", "math/big", "math/bits", "math/cmplx",
		"math/rand", "mime", "mime/multipart", "mime/quotedprintable",
		"net", "net/http", "net/http/cgi", "net/http/cookiejar", "net/http/fcgi",
		"net/http/httptest", "net/http/httptrace", "net/http/httputil", "net/http/pprof", "net/mail",
		"net/rpc", "net/rpc/jsonrpc", "net/smtp", "net/textproto", "net/url",
		"os", "os/exec", "os/signal", "os/user", "path",
		"path/filepath", "plugin", "reflect", "regexp", "regexp/syntax",
		"runtime", "runtime/cgo", "runtime/debug", "runtime/metrics", "runtime/msan",
		"runtime/pprof", "runtime/race", "runtime/trace", "sort", "strconv",
		"strings", "sync", "sync/atomic", "syscall", "syscall/js",
		"text/scanner", "text/tabwriter", "text/template", "text/template/parse", "time",
		"time/tzdata", "unicode", "unicode/utf16", "unicode/utf8", "unsafe",
	}
)

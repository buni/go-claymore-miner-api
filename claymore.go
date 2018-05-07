package claymore

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

const (
	// stats   = `{"id":0,"jsonrpc":"2.0","method":"miner_getstat1","psw":""}`
	reboot  = `{"id":0,"jsonrpc":"2.0","method":"miner_reboot","psw":""}`
	restart = `{"id":0,"jsonrpc":"2.0","method":"miner_restart","psw":""}`
)

//JSON input
type JSON struct {
	ID     int         `json:"id"`
	Error  interface{} `json:"error"`
	Result []string    `json:"result"`
}

//Stats struct
type Stats struct {
	Version    string    `json:"version"`
	Runtime    string    `json:"runtime"`
	SharesObj  []Shares  `json:"shares"`
	GPUMainObj []GPUMain `json:"gpumain"`
	SharesDObj []SharesD `json:"sharesd"`
	GPUDualObj []GPUDual `json:"gpudual"`
	TermalsObj []Termals `json:"Termals"`
	PoolsObj   []Pools   `json:"pools"`
	// PoolSharesObj []PoolShares `json:"poolshares"`
}

//Shares struct
type Shares struct {
	TotalHash int `json:"totalhash"`
	Accepted  int `json:"accepted"`
	Rejected  int `json:"rejected"`
}

//SharesD Dual mining shares struct
type SharesD struct {
	TotalHash int `json:"totalhash"`
	Accepted  int `json:"accepted"`
	Rejected  int `json:"rejected"`
}

//GPUMain main algo hashrate for each gpu
type GPUMain struct {
	Hashrate int `json:"hashrate"`
}

//GPUDual 2nd algo hashrate for each gpu
type GPUDual struct {
	Hashrate int `json:"hashrate"`
}

//Termals gpu termals & fan speed
type Termals struct {
	Temp int `json:"temp"`
	Fan  int `json:"fan"`
}

//Pools struct
type Pools struct {
	Pool string `json:"pool"`
}

// type PoolShares struct {
// 	Rejected     int `json:"rejected"`
// 	PoolSwtiches int `json:"poolswitches"`
// }

//GetConsole output
func GetConsole(url string) (sa string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	ht, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	html := string(ht)
	sa = sanitize(html)
	return sa, err
}

//Sanitize the console output
func sanitize(hts string) (res string) {
	doc, err := html.Parse(strings.NewReader(hts))
	if err != nil {
		log.Println(err)
	}
	removeScript(doc)
	buf := bytes.NewBuffer([]byte{})
	if err := html.Render(buf, doc); err != nil {
		log.Println(err)
	}
	htm := strings.TrimLeft(strings.TrimRight(buf.String(), "</html>"), "<html>")
	htm = strings.TrimLeft(strings.TrimRight(htm, "</body"), `body  bgcolor="#000000" style="font-family: monospace;>`)
	re := regexp.MustCompile("{.*} *")
	res = re.ReplaceAllString(htm, "")
	return
}

//removeScript tag
func removeScript(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "script" || n.Data == "head" {
		n.Parent.RemoveChild(n)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		removeScript(c)
	}
}

//GetStats get miner's stats
func GetStats(ip string, password string) (stat Stats, err error) {
	d := &net.Dialer{Timeout: time.Duration(5) * time.Millisecond}
	stats := `{"id":0,"jsonrpc":"2.0","method":"miner_getstat1","psw":"` + password + `"}`
	conn, err := d.Dial("tcp", ip)
	if err != nil {
		return
	}
	conn.Write([]byte(stats))
	// log.Printf("Send: %s", stats)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	// log.Printf("Receive: %s", buff[:n])
	defer conn.Close()
	stat = Normalize(buff[:n])
	return
}

//RestartMiner restart the miner(do not use it crashesh the miner more often then not)
func RestartMiner(ip string, password string) error {
	restart := `{"id":0,"jsonrpc":"2.0","method":"miner_restart","psw":"` + password + `"}`
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(err)
		return err
	}
	conn.Write([]byte(restart))
	log.Printf("Send: %s", restart)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
	defer conn.Close()
	return err
}

//RebootMiner reboots the miner(if there is no restart.sh/bat/bash it restarts the miner)
func RebootMiner(ip string, password string) error {
	reboot := `{"id":0,"jsonrpc":"2.0","method":"miner_reboot","psw":"` + password + `"}`

	conn, err := net.Dial("tcp", ip)
	if err != nil {
		log.Println(err)
		return err
	}
	conn.Write([]byte(reboot))
	log.Printf("Send: %s", reboot)
	buff := make([]byte, 1024)
	n, _ := conn.Read(buff)
	log.Printf("Receive: %s", buff[:n])
	defer conn.Close()
	return err
}

//["9.3 - ETH", "21",
//"182724;51;0", "30502;30457;30297;30481;30479;30505",
//"0;0;0", "off;off;off;off;off;off", "53;71;57;67;61;72;55;70;59;71;61;70",
//"eth-eu1.nanopool.org:9999", "0;0;0;0"]

//Normalize make the json useful
func Normalize(buff []byte) (stat Stats) {
	data := JSON{}
	json.Unmarshal(buff, &data)
	res := data.Result
	for i := 0; i < len(res); i++ {
		switch i {
		case 0:
			stat.Version = res[i]
		case 1:
			stat.Runtime = res[i]
			// ru, _ := strconv.Atoi(res[i])
			// d := ru / 60 / 24
			// h := ru / 60 % 24
			// m := ru % 60
			// fmt.Println(d, h, m)
		case 2:
			shares := castShares(res[i])
			stat.SharesObj = append(stat.SharesObj, shares)
		case 3:
			// gpumain := castGPUMain(res[i])
			stat.GPUMainObj = castGPUMain(res[i])
		case 4:
			shares, err := castSharesD(res[i])
			if err == nil {
				stat.SharesDObj = append(stat.SharesDObj, shares)
			}
		case 5:
			gpudual, err := castGPUDual(res[i])
			if err == nil {
				stat.GPUDualObj = gpudual
			}
		case 6:
			stat.TermalsObj = castTermals(res[i])
		case 7:
			stat.PoolsObj = castPools(res[i])
		case 8:
		}
	}
	return
}

//castShares
func castShares(data string) (shares Shares) {
	x := strings.Split(data, ";")
	h, _ := strconv.Atoi(x[0])
	a, _ := strconv.Atoi(x[1])
	r, _ := strconv.Atoi(x[2])
	shares.TotalHash = h
	shares.Accepted = a
	shares.Rejected = r
	return
}

//castGPUMain
func castGPUMain(data string) (gpumains []GPUMain) {
	var gpumain GPUMain
	x := strings.Split(data, ";")
	for i := 0; i < len(x); i++ {
		h, _ := strconv.Atoi(x[i])
		gpumain.Hashrate = h
		gpumains = append(gpumains, gpumain)
	}
	return
}

//castSharesD
func castSharesD(data string) (shares SharesD, err error) {

	x := strings.Split(data, ";")
	shares.TotalHash, err = strconv.Atoi(x[0])
	shares.Accepted, err = strconv.Atoi(x[1])
	shares.Rejected, err = strconv.Atoi(x[2])
	return shares, err
}

//castGPUDual
func castGPUDual(data string) (gpuduals []GPUDual, err error) {
	var gpudual GPUDual
	x := strings.Split(data, ";")
	for i := 0; i < len(x); i++ {
		if x[i] == "" || x[i] == "off" {
			return gpuduals, errors.New("an error")
		}
		h, _ := strconv.Atoi(x[i])
		gpudual.Hashrate = h
		gpuduals = append(gpuduals, gpudual)
	}
	return gpuduals, err
}

//castTermals
func castTermals(data string) (tss []Termals) {
	// data = `53;71;57;67;61;72;55;70;59;71;59;71;53;71;57;67;61;72;55;70;59;71;59;71`
	var ts Termals
	x := strings.Split(data, ";")
	// fmt.Println(x)
	if "" == data {
		return tss
	}
	for i := 0; i < len(x); i++ {
		if len(x) > i+1 {
			temp, _ := strconv.Atoi(x[i])
			ts.Temp = temp
			fan, _ := strconv.Atoi(x[i+1])
			ts.Fan = fan
			tss = append(tss, ts)
			i++
		}
	}
	return
}

//castPools
func castPools(data string) (pools []Pools) {
	var pool Pools
	x := strings.Split(data, ";")
	if len(x) > 0 {
		for i := 0; i < len(x); i++ {
			pool.Pool = x[i]
			pools = append(pools, pool)
		}
	} else {
		pool.Pool = data
		pools = append(pools, pool)
	}
	return
}

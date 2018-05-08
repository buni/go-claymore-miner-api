package claymore

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
)

const (
	getStats     = `{"id":0,"jsonrpc":"2.0","method":"miner_getstat1","psw":"password"}`
	rebootMiner  = `{"id":0,"jsonrpc":"2.0","method":"miner_reboot","psw":"password"}`
	restartMiner = `{"id":0,"jsonrpc":"2.0","method":"miner_restart","psw":"password"}`
	response     = `{"result": ["9.3 - ETH", "21", "182724;51;0", "30502;30457;30297;30481;30479;30505", "0;0;0", "off;off;off;off;off;off", "53;71;57;67;61;72;55;70;59;71;61;70", "eth-eu1.nanopool.org:9999", "0;0;0;0"]}`
	htmlresponse = `<html><body bgcolor="#000000" style="font-family: monospace;">
	{"result": ["10.6 - ETH", "0", "4840;0;0", "4840", "145205;0;0", "145205", "55;57", "eth-eu.dwarfpool.com:8008;dcr.coinmine.pl:2222", "0;0;0;0", "0", "0", "0", "0", "0", "0"]}<br><br><font color="#ffffff">
	</font><br><font color="#ffffff">����������������������������������������������������������������ͻ
	</font><br>&nbsp;<font color="#ffffff">�     Claymore's Dual ETH + DCR/SC/LBC/PASC GPU Miner v10.6      �
	</font><br><font color="#ffffff">����������������������������������������������������������������ͼ
	</font><br><font color="#ffffff">
	</font><br><font color="#ffffff">ETH: 1 pool is specified
	</font><br><font color="#ffffff">Main Ethereum pool is eth-eu.dwarfpool.com:8008
	</font><br><font color="#ffffff">DCR: 1 pool is specified
	</font><br><font color="#ffffff">Main Decred pool is dcr.coinmine.pl:2222
	</font><br><font color="#ffffff">OpenCL initializing...
	
	</font><br><font color="#ffffff">AMD Cards available: 1 
	</font><br><font color="#00ff00">GPU #0: Tahiti (AMD Radeon R9 200 / HD 7900 Series), 3072 MB available, 32 compute units (pci bus 1:0:0)
	</font><br><font color="#ffffff">GPU #0 recognized as Radeon 280X
	</font><br><font color="#ffffff">POOL/SOLO version
	</font><br><font color="#ffffff">GPU #0: algorithm ASM
	</font><br><font color="#ffffff">No NVIDIA CUDA GPUs detected.
	</font><br><font color="#00ff00">Total cards: 1 
	</font><br><font color="#ffffff">ETH: Stratum - connecting to 'eth-eu.dwarfpool.com' <46.105.68.41> port 8008
	</font><br><font color="#00ffff">DUAL MINING MODE ENABLED: ETHEREUM+DECRED
	</font><br><font color="#00ff00">ETH: Stratum - Connected (eth-eu.dwarfpool.com:8008)
	</font><br><font color="#ffffff">ETH: eth-proxy stratum mode
	</font><br><font color="#ffffff">Watchdog enabled
	</font><br><font color="#ffffff">Remote management (READ-ONLY MODE) is enabled on port 3333
	</font><br><font color="#ffffff">
	</font><br>&nbsp;<font color="#ffffff"> DCR: Stratum - connecting to 'dcr.coinmine.pl' <5.39.75.53> port 2222
	</font><br><font color="#00ff00">ETH: Authorized
	</font><br>&nbsp;<font color="#00ff00"> DCR: Stratum - Connected (dcr.coinmine.pl:2222)
	</font><br><font color="#ffffff">Setting DAG epoch #185...
	</font><br>&nbsp;<font color="#00ff00"> DCR: Authorized
	</font><br><font color="#ffffff">Setting DAG epoch #185 for GPU0
	</font><br><font color="#ffffff">Create GPU buffer for GPU0
	</font><br><font color="#ffffff">ETH: 05/06/18-00:28:37 - New job from eth-eu.dwarfpool.com:8008
	</font><br><font color="#00ffff">ETH - Total Speed: 0.000 Mh/s, Total Shares: 0, Rejected: 0, Time: 00:00
	</font><br><font color="#00ffff">ETH: GPU0 0.000 Mh/s
	</font><br>&nbsp;<font color="#ffff00"> DCR - Total Speed: 0.000 Mh/s, Total Shares: 0, Rejected: 0
	</font><br>&nbsp;<font color="#ffff00"> DCR: GPU0 0.000 Mh/s
	</font><br><font color="#ffffff">ETH: 05/06/18-00:28:41 - New job from eth-eu.dwarfpool.com:8008
	</font><br><font color="#00ffff">ETH - Total Speed: 0.000 Mh/s, Total Shares: 0, Rejected: 0, Time: 00:00
	</font><br><font color="#00ffff">ETH: GPU0 0.000 Mh/s
	</font><br>&nbsp;<font color="#ffff00"> DCR - Total Speed: 0.000 Mh/s, Total Shares: 0, Rejected: 0
	</font><br>&nbsp;<font color="#ffff00"> DCR: GPU0 0.000 Mh/s
	</font><br><font color="#ffffff">GPU0 DAG creation time - 11110 ms
	</font><br><font color="#ffffff">Setting DAG epoch #185 for GPU0 done
	</font><br><font color="#ff00ff">GPU0 t=53C fan=53%%
	</font><br><font color="#ffffff">ETH: 05/06/18-00:29:03 - New job from eth-eu.dwarfpool.com:8008
	</font><br><font color="#00ffff">ETH - Total Speed: 4.840 Mh/s, Total Shares: 0, Rejected: 0, Time: 00:00
	</font><br><font color="#00ffff">ETH: GPU0 4.840 Mh/s
	</font><br>&nbsp;<font color="#ffff00"> DCR - Total Speed: 145.188 Mh/s, Total Shares: 0, Rejected: 0
	</font><br>&nbsp;<font color="#ffff00"> DCR: GPU0 145.188 Mh/s
	</font><br>&nbsp;<font color="#ffffff"> DCR: 05/06/18-00:29:05 - New job from dcr.coinmine.pl:2222
	</font><br><font color="#ffffff">
	</font><br><font color="#00ff00">GPU #0: Tahiti (AMD Radeon R9 200 / HD 7900 Series), 3072 MB available, 32 compute units (pci bus 1:0:0)
	</font><br><font color="#00ffff">ETH - Total Speed: 4.840 Mh/s, Total Shares: 0, Rejected: 0, Time: 00:00
	</font><br><font color="#00ffff">ETH: GPU0 4.840 Mh/s
	</font><br>&nbsp;<font color="#ffff00"> DCR - Total Speed: 145.205 Mh/s, Total Shares: 0, Rejected: 0
	</font><br>&nbsp;<font color="#ffff00"> DCR: GPU0 145.205 Mh/s
	</font><br><font color="#00ffff">Incorrect ETH shares: none
	</font><br>&nbsp;<font color="#00ffff">1 minute average ETH total speed: 4.817 Mh/s
	</font><br><font color="#ffffff">Pool switches: ETH - 0, DCR - 0
	</font><br><font color="#ffffff">Current ETH share target: 0x0000000225c17d04 (diff: 2000MH), epoch 185(2.45GB)
	Current DCR share target: 0x0000000019998000 (diff: 42GH), block #236126
	</font><br><font color="#ff00ff">GPU0 t=55C fan=57%%
	</font><br><font color="#ffffff">
	</font><br></body></html>`
)
func init(){
        go StartHTTP()
}

var start = StartTCP()

func StartHTTP() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, htmlresponse)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func StartTCP() string {
	listen, err := net.Listen("tcp", "localhost:3306")
	if err != nil {
		log.Println("Error listening:", err)
	}
	go func() {
		for {
			conn, err := listen.Accept()
			if err != nil {
				log.Println("An error occured", err)
			}
			buf := make([]byte, 1024)
			_, err = conn.Read(buf)
			if err != nil {
				log.Println("An error occured:", err)
			}
			// bb := bytes.NewBuffer(buf).String()
			// fmt.Println(bb)
			conn.Write([]byte(response))
			conn.Close()
		}
	}()
	return ""
}

func TestGetStats(t *testing.T) {
	stat, err := GetStats("127.0.0.1:3306", "password")
	if err != nil {
		t.Error(err)
	}
	if stat.Version != "9.3 - ETH" {
		t.Error("Error with normalization")
	}

}
func TestRestartMiner(t *testing.T) {
	err := RestartMiner("127.0.0.1:3306", "password")
	if err != nil {
		t.Error(err)
	}

}
func TestRebootMiner(t *testing.T) {
	err := RebootMiner("127.0.0.1:3306", "password")
	if err != nil {
		t.Error(err)
	}
}
func TestNormalize(t *testing.T) {
	stats := Normalize([]byte(response))
	if stats.Version != "9.3 - ETH" {
		t.Error("Error with normalization")
	}
}
func TestGetConsole(t *testing.T) {
//	go StartHTTP()
//	time.Sleep(100 * time.Millisecond)
	output, err := GetConsole("http://localhost:8080")
	if err != nil || output == "" {
		t.Error("Error with GetConsole", err)
	}
}

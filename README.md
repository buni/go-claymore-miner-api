# Go based libary for dealing with Claymore Miner's API
## Usage
`go get github.com/buni/claymore-stats`
The package exposes five functions:
GetConsole(url,password) - Get miner's console output.
GetStats(ip,password) - Get miner's stats .
RestartMiner(ip,password) - Restart the miner(often crashes it)
RebootMiner(ip,password) - Reboots the miner(if there is no restart.sh/bat/bash it restarts the miner)
Normalize([]byte) - Normalize miner's response.
(Leave password blank if you have not set one yet.)
(The IP string should include both IP and port `GetStats("127.0.0.1:3306","")`)
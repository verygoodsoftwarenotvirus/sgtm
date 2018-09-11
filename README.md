# sgtm
hackathon screen reading project

TO BUILD:
* Clone repo into ~/go/src/github.com/verygoodsoftwarenotvirus/sgtm
* Run make revendor in terminal in that dir
* Run cli/main.go or playground/main.go

TO RUN CLI COMMANDS:
* "go run cmd/cli/main.go" in the terminal acts the same as "sgtm"
* The read function currently only accepts absolute paths
* For example, to run the read function run "go run cmd/cli/main.go read -f /example.go -p someFunction"
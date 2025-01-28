# Blocktime Node

A web server connecting to bitcoin core to display the node's blocktime in a browser.

## Test

Run Bitcoin Core 28.0 in regtest mode:

```bash
bitcoind -regtest
```

Create a wallet and generate some blocks:

```bash
bitcoin-cli -regtest createwallet "testwallet"
bitcoin-cli -regtest -generate 101
```

Run server:

```bash
go run cmd/server/main.go --rpc-cookie-file=$HOME/.bitcoin/regtest/.cookie --rpc-url=http://localhost:18443 --notify-socket=/tmp/blocktime-node.sock
```

Visit http://localhost:8080

## Build and install

```bash
make
sudo make install
blocktime-node --help
```

Create directory for notify socket file:

```bash
sudo mkdir /var/lib/blocktime-node
sudo chown bitcoin:bitcoin /var/lib/blocktime-node
```

Setup systemd service:

```bash
sudo cp config/systemd/blocktime-node.service /etc/systemd/system
# edit /etc/systemd/system/blocktime-node.service as required
sudo systemctl daemon-reload
sudo systemctl enable blocktime-node
sudo systemctl start blocktime-node
sudo systemctl status blocktime-node
```

Setup nginx proxy:

```bash
sudo cp config/nginx/blocktime-node /etc/nginx/sites-available
sudo ln -s /etc/nginx/sites-available/blocktime-node /etc/nginx/sites-enabled
# edit /etc/nginx/sites-enabled/blocktime-node as required
sudo nginx -t
sudo systemctl reload nginx
```

Setup block notifications with the `-blocknotify` flag (consider setting `blocknotify` in `bitcoin.conf` instead):

```bash
# In a production setup remove -regtest flag
bitcoind -regtest -blocknotify="blocktime-node-notify"
```

Load wallet and generate some blocks:

```bash
bitcoin-cli -regtest loadwallet "testwallet"
bitcoin-cli -regtest -generate 1
```

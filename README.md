# Blocktime Node

Blocktime Node is a web server that runs alongside bitcoin core to display the node's blocktime (technically block height) in a browser. It supports server-side HTML rendering to display the current blocktime when the web page is loaded or reloaded even if JavaScript is disabled in the browser and in addition supports server-sent events (SSE) for live updates if JavaScript is enabled.

The service is developed in go without any additional third-party libraries, making it easy to review and deploy.

## Test Setup

Run Bitcoin Core 28.0 in regtest mode:

```bash
bitcoind -regtest
```

Create a wallet and generate some blocks:

```bash
bitcoin-cli -regtest createwallet "testwallet"
bitcoin-cli -regtest -generate 42
```

Run server:

```bash
go run cmd/server/main.go --rpc-cookie-file=$HOME/.bitcoin/regtest/.cookie --rpc-url=http://localhost:18443 --notify-socket=/tmp/blocktime-node.sock
```

Open http://localhost:8080 in a browser.

While the `blocktime-node` server is running, it will only load the current blocktime from bitcoin core the first time the website is being loaded. It will only update the blocktime if the companion command `blocktime-node-notify` is being executed to trigger a blocktime update. This can be simulated by running the following command:

```bash
bitcoin-cli -regtest -generate 21
go run cmd/notify/main.go --notify-socket=/tmp/blocktime-node.sock
```

Use the `-blocknotify` setting of bitcoin core to automatically trigger blocktime updates via the `blocktime-node-notify` command:

```bash
make
sudo make install
bitcoind -regtest -blocknotify="blocktime-node-notify --notify-socket=/tmp/blocktime-node.sock"
bitcoin-cli -regtest loadwallet "testwallet"
bitcoin-cli -regtest -generate 1
```

The `blocktime-node` server will be notified via the notify-socket, connect to the bitcoin core via json-rpc to fetch the current blocktime, cache the blocktime for all future page loads and SSE connections, and send blocktime events to all connected SSE clients.

## Production Setup

Assumptions:

* git repository is cloned and is the current working directory
* [go](https://go.dev/doc/install) is installed
* bitcoin core is installed on a linux server
  * runs as user:group `bitcoin:bitcoin`
  * is managed by systemd as `bitcoind.service`
  * cookie file is generated at `/var/lib/bitcoind/.cookie`
  * config file is stored at `/etc/bitcoin/bitcoin.conf`
* nginx is installed and can be used as reverse proxy

Build and install `blocktime-node` and `blocktime-node-notify` commands:

```bash
make
sudo make install
# check configuration defaults and use commandline flags if assumptions do not apply
blocktime-node --help
blocktime-node-notify --help
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
sudo journalctl -u blocktime-node
```

Setup nginx proxy (and consider enabling HTTPS, instructions not included in the sample configuration):

```bash
sudo cp config/nginx/blocktime-node /etc/nginx/sites-available
sudo ln -s /etc/nginx/sites-available/blocktime-node /etc/nginx/sites-enabled
# edit /etc/nginx/sites-enabled/blocktime-node as required
sudo nginx -t
sudo systemctl reload nginx
```

Configure block notifications by editing `/etc/bitcoin/bitcoin.conf`:

```conf
blocknotify=blocktime-node-notify
```

Restart bitcoin core:

```bash
sudo systemctl restart bitcoind
sudo systemctl status bitcoind
```

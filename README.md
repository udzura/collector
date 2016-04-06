# collector

Collector of global IP and put them into DNS

## How to use

* Install binary

```
wget https://github.com/udzura/collector/releases/download/v0.2.0/collector_v0.2.0-linux-amd64.zip
unzip collector_v0.2.0-linux-amd64.zip
sudo cp collector_v0.2.0-linux-amd64 /usr/local/bin/collector
```

* collector is dependent on [consul](https://www.consul.io/), so create its cluster.
* Add the consul's check for each client instances, like:

```
{
  "service": {
    "id": "nginx",
    "name": "nginx",
    "tags": ["nginx"],
    "port": 80,
    "check":{
      "script": "/usr/local/bin/collector client --dev eth0 -- /usr/lib64/nagios/plugins/check_http -H localhost",
      "interval": "30s"
    }
  }
}
```

* Daemonize the consul's watch process in consumer instance, like:

```console
$ cat /root/.aws/credentials 
[default]
aws_access_key_id = AKIXXXXXX...
aws_secret_access_key = 4Jr...............
$ export SLACK_URL=https://your.slack.com/your/incoming/webhook-url
$ consul watch -type checks -service nginx -- \
      /usr/local/bin/collector watch --hosted-zone foo.example.com --domain front.foo.example.com
## Recommended to use systemc Unit file or like.
```

* Then, your consul automatically get Nginx health, and then update IPs on your Route53.

### Command options details

* `collector client` has some options:
  * Pass `--dev/-D` to specify device which has your global IP.
  * Default to `eth0`, right?

* `collector watch` is intended to run under `consul watch`
  * AWS credential info (`.aws/credential` or `AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY/AWS_REGION` environ) is required
  * If you set `SLACK_URL` environ, then changes are notified.

## Note

* This product is exoerimental and before alpha release.

## Contributing

* Usual GitHub way.

## License

* See LISENCE.

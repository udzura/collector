# collector

Collector of global IP and put them into DNS

[![Build Status](https://travis-ci.org/udzura/collector.svg)](https://travis-ci.org/udzura/collector)

## How to use

* Install binary

```bash
wget https://github.com/udzura/collector/releases/download/v0.3.3/collector_v0.3.3-linux-amd64.zip
unzip collector_v0.3.3-linux-amd64.zip
sudo cp collector_v0.3.3-linux-amd64 /usr/local/bin/collector
```

* collector is dependent on [consul](https://www.consul.io/), so create its cluster.
* Add the consul's check for each client instances, like:

```json
{
  "service": {
    "id": "nginx",
    "name": "nginx",
    "tags": ["nginx", "lb-a"],
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
$ consul watch -type service -service nginx -- \
      /usr/local/bin/collector watch --hosted-zone foo.example.com --domain front.foo.example.com
## Recommended to use systemc Unit file or like.
```

* Then, your consul automatically get Nginx health, and then update IPs on your Route53.

### Manage multi domain in one watch process

Create check with tagged `'lb-a'`:

```json
{
  "service": {
    "id": "nginx",
    "name": "nginx",
    "tags": ["nginx", "lb-a"],
    "port": 80,
    "check":{
      "script": "/usr/local/bin/collector client --dev eth0 -- /usr/lib64/nagios/plugins/check_http -H localhost",
      "interval": "30s"
    }
  }
}
```

At another host, reate check with tagged `'lb-b'`:

```json
{
  "service": {
    "id": "nginx",
    "name": "nginx",
    "tags": ["nginx", "lb-b"],
    "port": 80,
    "check":{
      "script": "/usr/local/bin/collector client --dev eth0 -- /usr/lib64/nagios/plugins/check_http -H localhost",
      "interval": "30s"
    }
  }
}
```

Then, pass multi `--domain` option with tag to consul watch:

```console
$ consul watch -type service -service nginx -- \
      /usr/local/bin/collector watch --hosted-zone foo.example.com \
      --domain front-a.foo.example.com:lb-a \
      --domain front-b.foo.example.com:lb-b
```

After this, check with `lb-a` effects domain `front-a.foo.example.com`, and `lb-b` effects `front-b.foo.example.com`.

### Command options details

* `collector client` has some options:
  * Pass `--dev/-D` to specify device which has your global IP.
  * Default to `eth0`, right?

* `collector watch` is intended to run under `consul watch`
  * AWS credential info (`.aws/credential` or `AWS_ACCESS_KEY_ID/AWS_SECRET_ACCESS_KEY/AWS_REGION` environ) is required
  * If you set `SLACK_URL` environ, then changes are notified.
  * Pass `--check-id` if your health check has CheckID which is a consul default `service:{{service_name}}`.

* `collector watch` respects environment variables:
  * Which are useful working with systemd unit file
  * `COLLECTOR_HOSTED_ZONE`, `COLLECTOR_DOMAIN` and `COLLECTOR_CHECK_ID`.
  * Note: `COLLECTOR_DOMAIN` should be splited with white space `" "`

## Note

* This product is experimental and before alpha release.

## Contributing

* Usual GitHub way.

## License

* See LISENCE.

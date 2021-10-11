# mackerel-plugin-nvme

it is alpha quality.

## usage

```
[plugin.metrics.nvme]
command = "nvme smart-log /dev/nvme0 --output-format=json | /usr/local/bin/mackerel-plugin-nvme"
```

need nvme-cli

## reference

https://www.ibm.com/docs/ja/power9?topic=devices-running-linux-smart-log-command


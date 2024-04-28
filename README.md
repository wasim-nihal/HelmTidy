# helm-dangling-tpls :  Find Unused Helm Chart Templates

This Go based tool scans Helm charts and reports unused named templates within those charts. It can also handle library charts (charts imported by other charts).

### Features

* Identifies unused named templates in Helm charts.
* Supports scanning library charts.
* Accepts absolute, relative, and remote chart paths (http/https).

### Usage

```
danglingTpls [flags]
```

**Flags:**

* `-p` string (required): path of the chart that has to be scanned for the unused tpls. It can be either relative/absolute/http url. If the value starts with `http://` or `https://`, then a http request is made to get the chart.
* `-dependentCharts` string (optional): a list comma saperated paths of dependent charts that make use of tpls from the current '-p' chart. Values can be either relative/absolute/http url. If the value starts with `http://` or `https://`, then a http request is made to get the chart.
* `-exitWithNonZero` bool (optional): flag indicating whether to exit with a non zero error code in the case if dangling tpls are found. (default: false)
* `-workerThreads` int (optional): Number of worker threads to use (default: 4).


**Building from Source:**

* Clone the repository
`git clone https://github.com/wasim-nihal/helm-dangling-tpls.git`

* Run the build command

```
cd helm-dangling-templates

make build
```

* The built binary can be found under `bin` directory

**Sample Usage:**

```
./danglingTpls -p=/path/to/my/chart
```
This scans the chart located at `/path/to/my/chart` for unused templates.

**Sample Usage for library/parent chart:**
```
./danglingTpls -p='/mnt/common-lib' --dependentCharts='/mnt/depChart1,/mnt/depChart2'
```

**Sample Usage with http url:**
```
./danglingTpls -p=http://<foo>/<bar>/<chart>.tgz
```

**Sample Output:**

```
$ ./danglingTpls -p=/mnt/c/helm-dangling-tpls/examples/charts/dummy-chart
2024/04/29 01:27:19 Scanning for tpl usages for the charts below:
2024/04/29 01:27:19 /mnt/c/helm-dangling-tpls/examples/charts/dummy-chart
+-----------------------+------------------------------------------------------------------------------------------+
| DANGLING TPL NAME     | FILE                                                                                     |
+-----------------------+------------------------------------------------------------------------------------------+
| dummy-chart.dummyTpl2 | /mnt/c/helm-dangling-tpls/examples/charts/dummy-chart/templates/_helpers.tpl             |
+-----------------------+------------------------------------------------------------------------------------------+
| dummy-chart.dummyTpl1 | /mnt/c/helm-dangling-tpls/examples/charts/dummy-chart/templates/_helpers.tpl             |
+-----------------------+------------------------------------------------------------------------------------------+
| dummy-chart.dummyTpl3 | /mnt/c/helm-dangling-tpls/examples/charts/dummy-chart/templates/anotherTpls/_example.tpl |
+-----------------------+------------------------------------------------------------------------------------------+
```
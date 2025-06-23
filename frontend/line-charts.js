// 定义初始化函数
function initializeLineCharts(charts) {
    fetch('/ops/ins-cache')
        .then(response => response.json())
        .then(data => {
            console.log('Received data:', data);
            const inspections = data.data;

            // 按主机名分类数据
            const hostData = {};
            inspections.forEach(record => {
                const hostname = record.hostname;

                if (!hostData[hostname]) {
                    hostData[hostname] = {
                        timestamps: [],
                        cpuUsage: [],
                        memoryUsed: [],
                        diskUsed: [],
                        networkSent: [],
                        networkRecv: [],
                        processesCount: []
                    };
                }

                const timestamp = new Date(record.timestamp);
                const formattedTime = timestamp.toLocaleString('zh-CN', {
                    year: 'numeric',
                    month: '2-digit',
                    day: '2-digit',
                    hour: '2-digit',
                    minute: '2-digit'
                });

                hostData[hostname].timestamps.push(formattedTime);
                hostData[hostname].cpuUsage.push(record.cpu?.usage || 0);
                hostData[hostname].memoryUsed.push(record.memory?.used || 0);
                hostData[hostname].diskUsed.push(record.disk[0]?.used || 0);
                hostData[hostname].networkSent.push(record.network[0]?.sent || 0);
                hostData[hostname].networkRecv.push(record.network[0]?.recv || 0);
                hostData[hostname].processesCount.push((record.processes || 0) + (record.zombieProcs || 0));
            });

            // 对每个主机的数据按时间排序
            for (const hostname in hostData) {
                const metrics = hostData[hostname];
                const sortedIndices = metrics.timestamps
                    .map((time, index) => ({ time, index }))
                    .sort((a, b) => new Date(a.time) - new Date(b.time))
                    .map(item => item.index);

                metrics.timestamps = sortedIndices.map(i => metrics.timestamps[i]);
                metrics.cpuUsage = sortedIndices.map(i => metrics.cpuUsage[i]);
                metrics.memoryUsed = sortedIndices.map(i => metrics.memoryUsed[i]);
                metrics.diskUsed = sortedIndices.map(i => metrics.diskUsed[i]);
                metrics.networkSent = sortedIndices.map(i => metrics.networkSent[i]);
                metrics.networkRecv = sortedIndices.map(i => metrics.networkRecv[i]);
                metrics.processesCount = sortedIndices.map(i => metrics.processesCount[i]);
            }

            // 绘制 CPU 使用率折线图
            charts['cpu-chart'].setOption({
                title: { text: 'CPU 使用率 (%)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'axis' },
                legend: { data: Object.keys(hostData) },
                xAxis: { type: 'category', data: hostData[Object.keys(hostData)[0]].timestamps },
                yAxis: { type: 'value' },
                series: Object.keys(hostData).map(hostname => ({
                    name: hostname,
                    type: 'line',
                    data: hostData[hostname].cpuUsage
                }))
            });

            // 绘制内存使用量折线图
            charts['memory-chart'].setOption({
                title: { text: '内存使用量 (MB)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'axis' },
                legend: { data: Object.keys(hostData) },
                xAxis: { type: 'category', data: hostData[Object.keys(hostData)[0]].timestamps },
                yAxis: { type: 'value' },
                series: Object.keys(hostData).map(hostname => ({
                    name: hostname,
                    type: 'line',
                    data: hostData[hostname].memoryUsed
                }))
            });

            // 绘制磁盘使用量折线图
            charts['disk-chart'].setOption({
                title: { text: '磁盘使用量 (MB)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'axis' },
                legend: { data: Object.keys(hostData) },
                xAxis: { type: 'category', data: hostData[Object.keys(hostData)[0]].timestamps },
                yAxis: { type: 'value' },
                series: Object.keys(hostData).map(hostname => ({
                    name: hostname,
                    type: 'line',
                    data: hostData[hostname].diskUsed
                }))
            });

            // 绘制网络流量折线图
            charts['network-chart'].setOption({
                title: { text: '网络流量 (MB)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'axis' },
                legend: { data: Object.keys(hostData) },
                xAxis: { type: 'category', data: hostData[Object.keys(hostData)[0]].timestamps },
                yAxis: { type: 'value' },
                series: Object.keys(hostData).map(hostname => ({
                    name: hostname,
                    type: 'line',
                    data: hostData[hostname].networkSent.map((sent, i) => sent + hostData[hostname].networkRecv[i])
                }))
            });

            // 绘制进程数量折线图
            charts['processes-chart'].setOption({
                title: { text: '进程数量 (个)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'axis' },
                legend: { data: Object.keys(hostData) },
                xAxis: { type: 'category', data: hostData[Object.keys(hostData)[0]].timestamps },
                yAxis: { type: 'value' },
                series: Object.keys(hostData).map(hostname => ({
                    name: hostname,
                    type: 'line',
                    data: hostData[hostname].processesCount
                }))
            });
        })
        .catch(error => {
            console.error('Error fetching data:', error);
            Object.values(charts).forEach(chart => chart.showLoading({ text: '数据加载失败', color: '#333' }));
        });
}
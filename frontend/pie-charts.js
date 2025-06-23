// 定义初始化函数
function initializePieCharts(charts) {
    fetch('/ops/latest-cache')
        .then(response => response.json())
        .then(data => {
            console.log('Received data:', data);
            const inspections = data.data;

            // CPU 使用率饼图
            charts['cpu-chart'].setOption({
                title: {
                    text: 'CPU 使用率 (%)',
                    left: 'center',
                    top: '5%',
                    textStyle: { color: '#333' }
                },
                tooltip: {
                    trigger: 'item',
                    formatter: '{a} <br/>{b}: {c}%'
                },
                series: [{
                    name: 'CPU 使用率',
                    type: 'pie',
                    radius: ['40%', '70%'],
                    center: ['50%', '60%'],
                    data: inspections.map(item => ({
                        name: item.hostname,
                        value: item.cpu?.usage || 0
                    })),
                    label: {
                        show: true,
                        formatter: '{b}: {c}%',
                        color: '#333'
                    },
                    emphasis: {
                        itemStyle: {
                            shadowBlur: 10,
                            shadowOffsetX: 0,
                            shadowColor: 'rgba(0, 0, 0, 0.5)'
                        }
                    }
                }]
            });

            // 内存使用量饼图
            charts['memory-chart'].setOption({
                title: { text: '内存使用量 (MB)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'item', formatter: '{a} <br/>{b}: {c} MB' },
                series: [{
                    name: '内存使用量',
                    type: 'pie',
                    radius: ['40%', '70%'],
                    center: ['50%', '60%'],
                    data: inspections.map(item => ({
                        name: item.hostname,
                        value: item.memory?.used || 0
                    })),
                    label: { show: true, formatter: '{b}: {c} MB', color: '#333' },
                    emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.5)' } }
                }]
            });

            // 磁盘使用量饼图
            charts['disk-chart'].setOption({
                title: { text: '磁盘使用量 (MB)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'item', formatter: '{a} <br/>{b}: {c} MB' },
                series: [{
                    name: '磁盘使用量',
                    type: 'pie',
                    radius: ['40%', '70%'],
                    center: ['50%', '60%'],
                    data: inspections.map(item => ({
                        name: item.hostname,
                        value: item.disk[0]?.used || 0
                    })),
                    label: { show: true, formatter: '{b}: {c} MB', color: '#333' },
                    emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.5)' } }
                }]
            });

            // 网络流量饼图
            charts['network-chart'].setOption({
                title: { text: '网络流量 (MB)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'item', formatter: '{a} <br/>{b}: {c} MB' },
                series: [{
                    name: '网络流量',
                    type: 'pie',
                    radius: ['40%', '70%'],
                    center: ['50%', '60%'],
                    data: inspections.map(item => ({
                        name: item.hostname,
                        value: ((item.network[0]?.sent || 0) + (item.network[0]?.recv || 0))
                    })),
                    label: { show: true, formatter: '{b}: {c} MB', color: '#333' },
                    emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.5)' } }
                }]
            });

            // 进程数量饼图
            charts['processes-chart'].setOption({
                title: { text: '进程数量 (个)', left: 'center', top: '5%', textStyle: { color: '#333' } },
                tooltip: { trigger: 'item', formatter: '{a} <br/>{b}: {c}' },
                series: [{
                    name: '进程数量',
                    type: 'pie',
                    radius: ['40%', '70%'],
                    center: ['50%', '60%'],
                    data: inspections.map(item => ({
                        name: item.hostname,
                        value: (item.processes || 0) + (item.zombieProcs || 0)
                    })),
                    label: { show: true, formatter: '{b}: {c}', color: '#333' },
                    emphasis: { itemStyle: { shadowBlur: 10, shadowOffsetX: 0, shadowColor: 'rgba(0, 0, 0, 0.5)' } }
                }]
            });
        })
        .catch(error => {
            console.error('Error fetching data:', error);
            Object.values(charts).forEach(chart => chart.showLoading({ text: '数据加载失败', color: '#333' }));
        });
}
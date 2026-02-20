import * as Humanize from 'humanize-plus'
import * as echarts from 'echarts/core'
import { PieChart, BarChart } from 'echarts/charts'
import { CanvasRenderer } from 'echarts/renderers'
import { LabelLayout } from 'echarts/features'

import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
} from 'echarts/components'

echarts.use([
  PieChart,
  BarChart,
  CanvasRenderer,
  LabelLayout,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent,
])

const isDark = () => document.documentElement.classList.contains('dark')

const getThemeTokens = () => {
  if (isDark()) {
    return {
      title: '#e2e8f0',
      subTitle: '#94a3b8',
      axisText: '#cbd5e1',
      splitLine: 'rgba(148, 163, 184, 0.2)',
      tooltipBg: 'rgba(15, 23, 42, 0.92)',
      palette: ['#38bdf8', '#2dd4bf', '#fbbf24', '#f87171', '#a78bfa', '#22c55e'],
      barIn: '#38bdf8',
      barOut: '#2dd4bf',
    }
  }

  return {
    title: '#0f172a',
    subTitle: '#475569',
    axisText: '#334155',
    splitLine: 'rgba(51, 65, 85, 0.14)',
    tooltipBg: 'rgba(15, 23, 42, 0.9)',
    palette: ['#0ea5e9', '#14b8a6', '#f59e0b', '#ef4444', '#8b5cf6', '#22c55e'],
    barIn: '#0ea5e9',
    barOut: '#14b8a6',
  }
}

const initChart = (elementId: string) => {
  const element = document.getElementById(elementId) as HTMLElement | null
  if (!element) {
    return null
  }

  const existing = echarts.getInstanceByDom(element)
  if (existing) {
    existing.dispose()
  }
  return echarts.init(element)
}

function DrawTrafficChart(
  elementId: string,
  trafficIn: number,
  trafficOut: number
) {
  const chart = initChart(elementId)
  if (!chart) {
    return
  }
  const tokens = getThemeTokens()

  chart.showLoading()

  chart.setOption({
    color: tokens.palette,
    title: {
      text: '网络流量',
      subtext: '今日统计',
      left: 'center',
      textStyle: { color: tokens.title, fontWeight: 700 },
      subtextStyle: { color: tokens.subTitle },
    },
    tooltip: {
      trigger: 'item',
      backgroundColor: tokens.tooltipBg,
      borderWidth: 0,
      textStyle: { color: '#f8fafc' },
      formatter: function (v: any) {
        return Humanize.fileSize(v.data.value) + ' (' + v.percent + '%)'
      },
    },
    legend: {
      orient: 'horizontal',
      bottom: 4,
      textStyle: { color: tokens.axisText },
      data: ['流量(入)', '流量(出)'],
    },
    series: [
      {
        type: 'pie',
        radius: ['45%', '70%'],
        center: ['50%', '50%'],
        padAngle: 2,
        itemStyle: {
          borderRadius: 8,
        },
        data: [
          {
            value: trafficIn,
            name: '流量(入)',
          },
          {
            value: trafficOut,
            name: '流量(出)',
          },
        ],
        emphasis: {
          itemStyle: {
            shadowBlur: 14,
            shadowOffsetX: 0,
            shadowColor: 'rgba(2, 6, 23, 0.35)',
          },
        },
      },
    ],
  })

  chart.hideLoading()
}

function DrawProxyChart(elementId: string, serverInfo: any) {
  const chart = initChart(elementId)
  if (!chart) {
    return
  }
  const tokens = getThemeTokens()

  chart.showLoading()

  const option = {
    color: tokens.palette,
    title: {
      text: '隧道分布',
      subtext: '按类型',
      left: 'center',
      textStyle: { color: tokens.title, fontWeight: 700 },
      subtextStyle: { color: tokens.subTitle },
    },
    tooltip: {
      trigger: 'item',
      backgroundColor: tokens.tooltipBg,
      borderWidth: 0,
      textStyle: { color: '#f8fafc' },
      formatter: function (v: any) {
        return String(v.data.value)
      },
    },
    legend: {
      orient: 'horizontal',
      bottom: 4,
      textStyle: { color: tokens.axisText },
      data: <string[]>[],
    },
    series: [
      {
        type: 'pie',
        radius: ['42%', '68%'],
        center: ['50%', '50%'],
        padAngle: 2,
        itemStyle: { borderRadius: 8 },
        data: <any[]>[],
      },
    ],
  }
  const proxyTypeCount = serverInfo?.proxyTypeCount || {}

  if (
    proxyTypeCount.tcp != null &&
    proxyTypeCount.tcp !== 0
  ) {
    option.series[0].data.push({
      value: proxyTypeCount.tcp,
      name: 'TCP',
    })
    option.legend.data.push('TCP')
  }
  if (
    proxyTypeCount.udp != null &&
    proxyTypeCount.udp !== 0
  ) {
    option.series[0].data.push({
      value: proxyTypeCount.udp,
      name: 'UDP',
    })
    option.legend.data.push('UDP')
  }
  if (
    proxyTypeCount.http != null &&
    proxyTypeCount.http !== 0
  ) {
    option.series[0].data.push({
      value: proxyTypeCount.http,
      name: 'HTTP',
    })
    option.legend.data.push('HTTP')
  }
  if (
    proxyTypeCount.https != null &&
    proxyTypeCount.https !== 0
  ) {
    option.series[0].data.push({
      value: proxyTypeCount.https,
      name: 'HTTPS',
    })
    option.legend.data.push('HTTPS')
  }
  if (
    proxyTypeCount.stcp != null &&
    proxyTypeCount.stcp !== 0
  ) {
    option.series[0].data.push({
      value: proxyTypeCount.stcp,
      name: 'STCP',
    })
    option.legend.data.push('STCP')
  }
  if (
    proxyTypeCount.sudp != null &&
    proxyTypeCount.sudp !== 0
  ) {
    option.series[0].data.push({
      value: proxyTypeCount.sudp,
      name: 'SUDP',
    })
    option.legend.data.push('SUDP')
  }
  if (
    proxyTypeCount.xtcp != null &&
    proxyTypeCount.xtcp !== 0
  ) {
    option.series[0].data.push({
      value: proxyTypeCount.xtcp,
      name: 'XTCP',
    })
    option.legend.data.push('XTCP')
  }

  chart.setOption(option)
  chart.hideLoading()
}

function DrawProxyTrafficChart(
  elementId: string,
  trafficInArr: number[],
  trafficOutArr: number[]
) {
  const chart = initChart(elementId)
  if (!chart) {
    return
  }
  const tokens = getThemeTokens()

  chart.showLoading()

  trafficInArr = trafficInArr.reverse()
  trafficOutArr = trafficOutArr.reverse()
  let now = new Date()
  now = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 6)
  const dates: Array<string> = []
  for (let i = 0; i < 7; i++) {
    dates.push(
      now.getFullYear() + '-' + (now.getMonth() + 1) + '-' + now.getDate()
    )
    now = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1)
  }

  chart.setOption({
    color: [tokens.barIn, tokens.barOut],
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow',
      },
      backgroundColor: tokens.tooltipBg,
      borderWidth: 0,
      textStyle: { color: '#f8fafc' },
      formatter: function (data: any) {
        let html = ''
        if (data.length > 0) {
          html += data[0].name + '<br/>'
        }
        for (const v of data) {
          const colorEl =
            '<span style="display:inline-block;margin-right:5px;' +
            'border-radius:10px;width:9px;height:9px;background-color:' +
            v.color +
            '"></span>'
          html +=
            colorEl + v.seriesName + ': ' + Humanize.fileSize(v.value) + '<br/>'
        }
        return html
      },
    },
    legend: {
      data: ['流量(入)', '流量(出)'],
      textStyle: { color: tokens.axisText },
    },
    grid: {
      left: '3%',
      right: '3%',
      bottom: '3%',
      top: '16%',
      containLabel: true,
    },
    xAxis: [
      {
        type: 'category',
        data: dates,
        axisLabel: {
          color: tokens.axisText,
        },
        axisLine: {
          lineStyle: {
            color: tokens.splitLine,
          },
        },
      },
    ],
    yAxis: [
      {
        type: 'value',
        axisLabel: {
          color: tokens.axisText,
          formatter: function (value: number) {
            return Humanize.fileSize(value)
          },
        },
        splitLine: {
          lineStyle: {
            color: tokens.splitLine,
          },
        },
      },
    ],
    series: [
      {
        name: '流量(入)',
        type: 'bar',
        barWidth: '34%',
        itemStyle: {
          borderRadius: [8, 8, 0, 0],
        },
        data: trafficInArr,
      },
      {
        name: '流量(出)',
        type: 'bar',
        barWidth: '34%',
        itemStyle: {
          borderRadius: [8, 8, 0, 0],
        },
        data: trafficOutArr,
      },
    ],
  })

  chart.hideLoading()
}

export { DrawTrafficChart, DrawProxyChart, DrawProxyTrafficChart }

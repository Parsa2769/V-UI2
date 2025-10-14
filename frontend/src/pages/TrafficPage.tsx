import { useState, useEffect } from 'react'
import { TrendingUp, TrendingDown, Activity, Calendar, Download, Upload } from 'lucide-react'
import { api } from '@/lib/api'
import { toast } from 'sonner'

interface TrafficData {
  date: string
  upload: number
  download: number
  total: number
}

interface TopClient {
  email: string
  total_traffic: number
  upload: number
  download: number
}

export default function TrafficPage() {
  const [loading, setLoading] = useState(true)
  const [timeRange, setTimeRange] = useState<'24h' | '7d' | '30d' | '90d'>('7d')
  const [trafficData, setTrafficData] = useState<TrafficData[]>([])
  const [topClients, setTopClients] = useState<TopClient[]>([])
  const [stats, setStats] = useState({
    total_upload: 0,
    total_download: 0,
    total: 0,
    avg_daily: 0,
  })

  useEffect(() => {
    fetchTrafficData()
  }, [timeRange])

  const fetchTrafficData = async () => {
    try {
      setLoading(true)
      const [trafficRes, topClientsRes, statsRes] = await Promise.all([
        api.get(`/traffic/history?range=${timeRange}`),
        api.get('/traffic/top-clients?limit=10'),
        api.get(`/traffic/stats?range=${timeRange}`),
      ])
      
      setTrafficData(trafficRes.data.data || [])
      setTopClients(topClientsRes.data.data || [])
      setStats(statsRes.data || {})
    } catch (error) {
      toast.error('Failed to fetch traffic data')
      console.error(error)
    } finally {
      setLoading(false)
    }
  }

  const formatBytes = (bytes: number) => {
    if (bytes === 0) return '0 B'
    const k = 1024
    const sizes = ['B', 'KB', 'MB', 'GB', 'TB']
    const i = Math.floor(Math.log(bytes) / Math.log(k))
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
  }

  const getMaxTraffic = () => {
    return Math.max(...trafficData.map(d => Math.max(d.upload, d.download)))
  }

  const getBarHeight = (value: number) => {
    const max = getMaxTraffic()
    return max > 0 ? (value / max) * 100 : 0
  }

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
      </div>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold">Traffic Statistics</h1>
          <p className="text-muted-foreground mt-1">
            Monitor bandwidth usage and traffic patterns
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Calendar className="w-4 h-4 text-muted-foreground" />
          <select
            value={timeRange}
            onChange={(e) => setTimeRange(e.target.value as any)}
            className="px-3 py-2 border border-border rounded-md bg-background"
          >
            <option value="24h">Last 24 Hours</option>
            <option value="7d">Last 7 Days</option>
            <option value="30d">Last 30 Days</option>
            <option value="90d">Last 90 Days</option>
          </select>
        </div>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">Total Upload</div>
            <Upload className="w-4 h-4 text-blue-600" />
          </div>
          <div className="text-2xl font-bold mt-1">{formatBytes(stats.total_upload)}</div>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">Total Download</div>
            <Download className="w-4 h-4 text-green-600" />
          </div>
          <div className="text-2xl font-bold mt-1">{formatBytes(stats.total_download)}</div>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">Total Traffic</div>
            <Activity className="w-4 h-4 text-purple-600" />
          </div>
          <div className="text-2xl font-bold mt-1">{formatBytes(stats.total)}</div>
        </div>
        <div className="bg-card border border-border rounded-lg p-4">
          <div className="flex items-center justify-between">
            <div className="text-sm text-muted-foreground">Avg Daily</div>
            <TrendingUp className="w-4 h-4 text-orange-600" />
          </div>
          <div className="text-2xl font-bold mt-1">{formatBytes(stats.avg_daily)}</div>
        </div>
      </div>

      {/* Traffic Chart */}
      <div className="bg-card border border-border rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-6">Traffic Over Time</h2>
        
        {/* Simple Bar Chart */}
        <div className="space-y-2">
          <div className="flex justify-between text-xs text-muted-foreground mb-2">
            <span>Date</span>
            <div className="flex gap-4">
              <span className="flex items-center gap-1">
                <div className="w-3 h-3 bg-blue-500 rounded"></div>
                Upload
              </span>
              <span className="flex items-center gap-1">
                <div className="w-3 h-3 bg-green-500 rounded"></div>
                Download
              </span>
            </div>
          </div>
          
          <div className="relative h-64">
            {trafficData.length > 0 ? (
              <div className="flex items-end justify-between h-full gap-1">
                {trafficData.map((data, index) => (
                  <div key={index} className="flex-1 flex flex-col items-center justify-end h-full group">
                    <div className="flex flex-col items-center justify-end w-full h-full gap-1">
                      {/* Upload Bar */}
                      <div 
                        className="w-full bg-blue-500 rounded-t hover:bg-blue-600 transition cursor-pointer relative"
                        style={{ height: `${getBarHeight(data.upload)}%` }}
                        title={`Upload: ${formatBytes(data.upload)}`}
                      >
                        {getBarHeight(data.upload) > 10 && (
                          <div className="absolute inset-0 flex items-center justify-center text-white text-xs opacity-0 group-hover:opacity-100 transition">
                            {formatBytes(data.upload)}
                          </div>
                        )}
                      </div>
                      {/* Download Bar */}
                      <div 
                        className="w-full bg-green-500 rounded-t hover:bg-green-600 transition cursor-pointer relative"
                        style={{ height: `${getBarHeight(data.download)}%` }}
                        title={`Download: ${formatBytes(data.download)}`}
                      >
                        {getBarHeight(data.download) > 10 && (
                          <div className="absolute inset-0 flex items-center justify-center text-white text-xs opacity-0 group-hover:opacity-100 transition">
                            {formatBytes(data.download)}
                          </div>
                        )}
                      </div>
                    </div>
                    <div className="text-xs text-muted-foreground mt-2 transform -rotate-45 origin-top-left">
                      {new Date(data.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' })}
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <div className="flex items-center justify-center h-full text-muted-foreground">
                No traffic data available
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Top Clients */}
      <div className="bg-card border border-border rounded-lg p-6">
        <h2 className="text-xl font-semibold mb-4">Top Clients by Traffic</h2>
        
        <div className="space-y-3">
          {topClients.map((client, index) => (
            <div key={index} className="flex items-center gap-4 p-3 hover:bg-accent rounded-lg transition">
              <div className="flex-shrink-0 w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center font-semibold text-primary">
                {index + 1}
              </div>
              <div className="flex-1 min-w-0">
                <div className="font-medium truncate">{client.email}</div>
                <div className="text-sm text-muted-foreground">
                  ↑ {formatBytes(client.upload)} · ↓ {formatBytes(client.download)}
                </div>
              </div>
              <div className="text-right">
                <div className="font-semibold">{formatBytes(client.total_traffic)}</div>
                <div className="text-xs text-muted-foreground">Total</div>
              </div>
            </div>
          ))}
          
          {topClients.length === 0 && (
            <div className="text-center py-8 text-muted-foreground">
              No client data available
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

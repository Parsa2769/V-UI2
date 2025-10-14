import { useQuery } from '@tanstack/react-query'
import { Users, Shield, Server, Activity } from 'lucide-react'
import api from '@/lib/api'

interface MetricsSummary {
  users: { total: number }
  clients: { total: number; active: number }
  nodes: { total: number; online: number }
}

export default function DashboardPage() {
  const { data: metrics } = useQuery<MetricsSummary>({
    queryKey: ['metrics-summary'],
    queryFn: async () => {
      const response = await api.get('/metrics/summary')
      return response.data
    },
    refetchInterval: 5000,
  })

  const cards = [
    {
      title: 'Total Users',
      value: metrics?.users.total || 0,
      icon: Users,
      color: 'blue',
    },
    {
      title: 'Active Clients',
      value: `${metrics?.clients.active || 0} / ${metrics?.clients.total || 0}`,
      icon: Shield,
      color: 'green',
    },
    {
      title: 'Online Nodes',
      value: `${metrics?.nodes.online || 0} / ${metrics?.nodes.total || 0}`,
      icon: Server,
      color: 'purple',
    },
    {
      title: 'Traffic',
      value: '0 GB',
      icon: Activity,
      color: 'orange',
    },
  ]

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
        <p className="text-gray-600 dark:text-gray-400 mt-1">
          Overview of your Xray infrastructure
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {cards.map((card) => (
          <div
            key={card.title}
            className="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700"
          >
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-gray-600 dark:text-gray-400 mb-1">
                  {card.title}
                </p>
                <p className="text-3xl font-bold text-gray-900 dark:text-white">
                  {card.value}
                </p>
              </div>
              <div className={`p-3 bg-${card.color}-100 dark:bg-${card.color}-900/20 rounded-lg`}>
                <card.icon className={`w-6 h-6 text-${card.color}-600 dark:text-${card.color}-400`} />
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
            Recent Activity
          </h2>
          <p className="text-gray-600 dark:text-gray-400">
            No recent activity to display
          </p>
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-lg shadow p-6 border border-gray-200 dark:border-gray-700">
          <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-4">
            System Status
          </h2>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">API Status</span>
              <span className="px-2 py-1 bg-green-100 dark:bg-green-900/20 text-green-600 dark:text-green-400 text-sm rounded">
                Healthy
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">Database</span>
              <span className="px-2 py-1 bg-green-100 dark:bg-green-900/20 text-green-600 dark:text-green-400 text-sm rounded">
                Connected
              </span>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-gray-600 dark:text-gray-400">Xray Core</span>
              <span className="px-2 py-1 bg-green-100 dark:bg-green-900/20 text-green-600 dark:text-green-400 text-sm rounded">
                Running
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

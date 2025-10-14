import { NavLink } from 'react-router-dom'
import { 
  LayoutDashboard, 
  Users, 
  Server, 
  Activity, 
  FileText, 
  Settings,
  Shield,
  Inbox,
  Lock,
  Send,
  Database,
  Layers
} from 'lucide-react'
import { cn } from '@/lib/utils'

const navigation = [
  { name: 'Dashboard', href: '/', icon: LayoutDashboard },
  { name: 'Users', href: '/users', icon: Users },
  { name: 'Clients', href: '/clients', icon: Shield },
  { name: 'Inbounds', href: '/inbounds', icon: Inbox },
  { name: 'Nodes', href: '/nodes', icon: Server },
  { name: 'Traffic', href: '/traffic', icon: Activity },
  { name: 'Certificates', href: '/certificates', icon: Lock },
  { name: 'Telegram', href: '/telegram', icon: Send },
  { name: 'Backup', href: '/backup', icon: Database },
  { name: 'Templates', href: '/templates', icon: Layers },
  { name: 'Audit Logs', href: '/audit', icon: FileText },
  { name: 'Settings', href: '/settings', icon: Settings },
]

export default function Sidebar() {
  return (
    <div className="w-64 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 flex flex-col">
      <div className="p-6 border-b border-gray-200 dark:border-gray-700">
        <div className="flex items-center space-x-3">
          <Shield className="w-8 h-8 text-blue-600 dark:text-blue-400" />
          <div>
            <h1 className="text-xl font-bold text-gray-900 dark:text-white">V-UI</h1>
            <p className="text-xs text-gray-500 dark:text-gray-400">v2.0.0</p>
          </div>
        </div>
      </div>

      <nav className="flex-1 p-4 space-y-1 overflow-y-auto">
        {navigation.map((item) => (
          <NavLink
            key={item.name}
            to={item.href}
            end={item.href === '/'}
            className={({ isActive }) =>
              cn(
                'flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors',
                isActive
                  ? 'bg-blue-50 dark:bg-blue-900/20 text-blue-600 dark:text-blue-400'
                  : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
              )
            }
          >
            {({ isActive }) => (
              <>
                <item.icon className={cn('w-5 h-5', isActive && 'text-blue-600 dark:text-blue-400')} />
                <span className="font-medium">{item.name}</span>
              </>
            )}
          </NavLink>
        ))}
      </nav>

      <div className="p-4 border-t border-gray-200 dark:border-gray-700">
        <p className="text-xs text-center text-gray-500 dark:text-gray-400">
          v2.0.0 • GPL-3.0
        </p>
      </div>
    </div>
  )
}

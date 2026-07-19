import React from 'react';
import { BrowserRouter, Routes, Route, NavLink } from 'react-router-dom';
import {
  LayoutDashboard, AlertTriangle, Shield,
  GitBranch, Crosshair, CheckSquare, Settings, Plug
} from 'lucide-react';

import OverviewPage      from './pages/OverviewPage';
import DetectionsPage    from './pages/DetectionsPage';
import { PoliciesPage }  from './pages/PoliciesPage';
import { TraceExplorerPage } from './pages/TraceExplorerPage';
import RedTeamPage       from './pages/RedTeamPage';
import { ProvidersPage } from './pages/ProvidersPage';
import CompliancePage    from './pages/CompliancePage';
import SettingsPage      from './pages/SettingsPage';

const NAV = [
  { to: '/',           icon: LayoutDashboard, label: 'Overview'       },
  { to: '/detections', icon: AlertTriangle,   label: 'Detections'     },
  { to: '/policies',   icon: Shield,          label: 'Policies'       },
  { to: '/traces',     icon: GitBranch,       label: 'Trace Explorer' },
  { to: '/redteam',    icon: Crosshair,       label: 'Red Team'       },
  { to: '/providers',  icon: Plug,            label: 'Providers'      },
  { to: '/compliance', icon: CheckSquare,     label: 'Compliance'     },
  { to: '/settings',   icon: Settings,        label: 'Settings'       },
];

const App = () => (
    <div className="flex h-screen bg-[#0a0a0f] text-gray-200" style={{ fontFamily: "'Inter', sans-serif" }}>

      {/* ── Sidebar ── */}
      <aside className="w-64 bg-[#12121a] border-r border-gray-800 flex flex-col flex-shrink-0">

        {/* Logo */}
        <div className="p-5 flex items-center gap-3">
          <Shield className="text-blue-500 w-8 h-8" />
          <span className="text-xl font-bold text-white tracking-widest">AEGIS</span>
        </div>

        {/* Nav */}
        <nav className="flex-1 px-3 py-4 space-y-1">
          {NAV.map(({ to, icon: Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              end={to === '/'}
              className={({ isActive }) =>
                `flex items-center gap-3 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${
                  isActive
                    ? 'bg-blue-500/10 text-blue-400 border border-blue-500/20'
                    : 'text-gray-400 hover:text-gray-200 hover:bg-gray-800/60'
                }`
              }
            >
              <Icon className="w-4 h-4 flex-shrink-0" />
              <span>{label}</span>
            </NavLink>
          ))}
        </nav>

        {/* Footer */}
        <div className="p-4 border-t border-gray-800">
          <div className="flex items-center gap-2">
            <span className="w-2 h-2 rounded-full bg-green-400 animate-pulse" />
            <span className="text-xs text-gray-500">v0.8.0 — All systems nominal</span>
          </div>
        </div>
      </aside>

      {/* ── Main content ── */}
      <main className="flex-1 overflow-auto">
        <Routes>
          <Route path="/"           element={<OverviewPage />}      />
          <Route path="/detections" element={<DetectionsPage />}    />
          <Route path="/policies"   element={<PoliciesPage />}      />
          <Route path="/traces"     element={<TraceExplorerPage />} />
          <Route path="/redteam"    element={<RedTeamPage />}       />
          <Route path="/providers"  element={<ProvidersPage />}     />
          <Route path="/compliance" element={<CompliancePage />}    />
          <Route path="/settings"   element={<SettingsPage />}      />
        </Routes>
      </main>
    </div>
);

export default App;

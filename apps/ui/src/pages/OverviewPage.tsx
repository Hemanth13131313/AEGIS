import React from 'react';
import { useDetectionSummary, useDetections } from '../hooks/useDetections';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer } from 'recharts';
import { formatDistanceToNow } from 'date-fns';

const mockChartData = [
  { day: 'Mon', detections: 4 },
  { day: 'Tue', detections: 7 },
  { day: 'Wed', detections: 3 },
  { day: 'Thu', detections: 12 },
  { day: 'Fri', detections: 8 },
  { day: 'Sat', detections: 2 },
  { day: 'Sun', detections: 11 },
];

export default function OverviewPage() {
  const { data: summary, isLoading: isSummaryLoading } = useDetectionSummary();
  const { data: recent, isLoading: isRecentLoading } = useDetections({ limit: 5 });

  return (
    <div className="p-6 space-y-6 max-w-7xl mx-auto">
      <h1 className="text-2xl font-bold text-gray-100">Overview</h1>
      
      {/* Stat Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        {[
          { label: 'Total Detections (Today)', value: summary?.total_today || 0 },
          { label: 'Active Policies', value: summary?.active_policies || 0 },
          { label: 'Protected Endpoints', value: summary?.endpoints || 0 },
          { label: 'Red Team Coverage', value: `${summary?.redteam_coverage || 0}%` },
        ].map((stat, i) => (
          <div key={i} className="bg-gray-800 border border-gray-700 rounded-xl p-5 shadow-sm">
            <div className="text-sm text-gray-400 mb-1">{stat.label}</div>
            <div className="text-3xl font-bold text-gray-100">
              {isSummaryLoading ? '...' : stat.value}
            </div>
          </div>
        ))}
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Chart */}
        <div className="lg:col-span-2 bg-gray-800 border border-gray-700 rounded-xl p-5 shadow-sm">
          <h2 className="text-lg font-semibold text-gray-200 mb-4">Detection Trend (7 Days)</h2>
          <div className="h-64">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={mockChartData}>
                <defs>
                  <linearGradient id="colorDet" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#6366f1" stopOpacity={0.3}/>
                    <stop offset="95%" stopColor="#6366f1" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <XAxis dataKey="day" stroke="#4b5563" tick={{fill: '#9ca3af'}} />
                <YAxis stroke="#4b5563" tick={{fill: '#9ca3af'}} />
                <Tooltip 
                  contentStyle={{ backgroundColor: '#1f2937', border: '1px solid #374151', color: '#f3f4f6' }}
                />
                <Area type="monotone" dataKey="detections" stroke="#6366f1" fillOpacity={1} fill="url(#colorDet)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Recent Detections List */}
        <div className="bg-gray-800 border border-gray-700 rounded-xl p-5 shadow-sm">
          <div className="flex justify-between items-center mb-4">
            <h2 className="text-lg font-semibold text-gray-200">Recent Activity</h2>
            <a href="/detections" className="text-indigo-400 hover:text-indigo-300 text-sm">View All</a>
          </div>
          
          <div className="space-y-4">
            {isRecentLoading ? (
               <div className="animate-pulse bg-gray-700 h-32 rounded w-full"></div>
            ) : (
              recent?.detections.slice(0, 5).map(d => (
                <div key={d.id} className="flex justify-between items-start pb-4 border-b border-gray-700 last:border-0">
                  <div>
                    <div className="flex items-center space-x-2 mb-1">
                      <span className={`w-2 h-2 rounded-full ${d.severity === 'critical' ? 'bg-red-500' : d.severity === 'high' ? 'bg-orange-500' : d.severity === 'medium' ? 'bg-yellow-500' : d.severity === 'low' ? 'bg-blue-500' : 'bg-green-500'}`}></span>
                      <span className="text-sm font-medium text-gray-200">{d.category}</span>
                    </div>
                    <div className="text-xs text-gray-500 font-mono">
                      {d.session_id.substring(0, 8)} • {d.owasp_category}
                    </div>
                  </div>
                  <div className="text-right">
                    <span className="text-xs text-gray-400 block mb-1">
                      {formatDistanceToNow(new Date(d.timestamp), { addSuffix: true })}
                    </span>
                    <span className="text-xs uppercase px-1.5 py-0.5 bg-gray-700 text-gray-300 rounded border border-gray-600">
                      {d.action_taken}
                    </span>
                  </div>
                </div>
              ))
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

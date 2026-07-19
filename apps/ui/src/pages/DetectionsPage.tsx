import React, { useState } from 'react';
import { useDetections, Detection } from '../hooks/useDetections';
import { formatDistanceToNow } from 'date-fns';

const SeverityBadge = ({ severity }: { severity: string }) => {
  const colors: Record<string, string> = {
    critical: 'bg-red-900 text-red-100',
    high: 'bg-orange-900 text-orange-100',
    medium: 'bg-yellow-900 text-yellow-100',
    low: 'bg-blue-900 text-blue-100',
    safe: 'bg-green-900 text-green-100',
  };
  return (
    <span className={`px-2 py-1 text-xs rounded-full font-medium ${colors[severity] || colors.safe}`}>
      {severity.toUpperCase()}
    </span>
  );
};

const ActionBadge = ({ action }: { action: string }) => {
  const colors: Record<string, string> = {
    block: 'border-red-500 text-red-400',
    redact: 'border-yellow-500 text-yellow-400',
    tag: 'border-blue-500 text-blue-400',
    allow: 'border-green-500 text-green-400',
  };
  return (
    <span className={`px-2 py-1 text-xs border rounded font-medium bg-gray-800 ${colors[action] || colors.allow}`}>
      {action.toUpperCase()}
    </span>
  );
};

export default function DetectionsPage() {
  const [live, setLive] = useState(false);
  const [severityFilter, setSeverityFilter] = useState('all');
  
  const { data, isLoading } = useDetections({ live, severity: severityFilter === 'all' ? undefined : severityFilter });

  return (
    <div className="p-6 space-y-6 max-w-7xl mx-auto">
      <div className="flex justify-between items-center">
        <h1 className="text-2xl font-bold text-gray-100">Detections</h1>
        <div className="flex items-center space-x-2">
          <label className="flex items-center space-x-2 text-sm text-gray-300">
            <input 
              type="checkbox" 
              checked={live} 
              onChange={e => setLive(e.target.checked)}
              className="rounded bg-gray-800 border-gray-600 text-indigo-500 focus:ring-indigo-500"
            />
            <span>Live Updates (5s)</span>
          </label>
        </div>
      </div>

      <div className="flex space-x-4 mb-4">
        <select 
          value={severityFilter}
          onChange={e => setSeverityFilter(e.target.value)}
          className="bg-gray-800 border border-gray-700 text-gray-200 rounded p-2 text-sm"
        >
          <option value="all">All Severities</option>
          <option value="critical">Critical</option>
          <option value="high">High</option>
          <option value="medium">Medium</option>
          <option value="low">Low</option>
          <option value="safe">Safe</option>
        </select>
        {/* Mock other filters */}
        <input 
          type="text" 
          placeholder="Search Session ID..." 
          className="bg-gray-800 border border-gray-700 text-gray-200 rounded p-2 text-sm w-64"
        />
      </div>

      {isLoading ? (
        <div className="space-y-4">
          {[1,2,3].map(i => (
            <div key={i} className="animate-pulse bg-gray-800 h-16 rounded-lg w-full"></div>
          ))}
        </div>
      ) : data?.detections.length === 0 ? (
        <div className="text-center py-12 bg-gray-800/50 rounded-lg border border-gray-700">
          <div className="text-4xl mb-4">🛡️</div>
          <p className="text-gray-400">No detections in this time range. All clear.</p>
        </div>
      ) : (
        <div className="bg-gray-800 border border-gray-700 rounded-lg overflow-hidden">
          <table className="w-full text-left text-sm text-gray-300">
            <thead className="bg-gray-900/50 text-gray-400 text-xs uppercase border-b border-gray-700">
              <tr>
                <th className="px-4 py-3">Severity</th>
                <th className="px-4 py-3">Category</th>
                <th className="px-4 py-3">OWASP / ATLAS</th>
                <th className="px-4 py-3">Action</th>
                <th className="px-4 py-3">Session ID</th>
                <th className="px-4 py-3">Time</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-700">
              {data?.detections.map((d: Detection) => (
                <tr key={d.id} className="hover:bg-gray-700/50 cursor-pointer transition-colors">
                  <td className="px-4 py-3"><SeverityBadge severity={d.severity} /></td>
                  <td className="px-4 py-3 font-medium text-gray-200">{d.category}</td>
                  <td className="px-4 py-3 text-xs">
                    <span className="text-indigo-400 mr-2">{d.owasp_category}</span>
                    <span className="text-gray-500">{d.atlas_technique}</span>
                  </td>
                  <td className="px-4 py-3"><ActionBadge action={d.action_taken} /></td>
                  <td className="px-4 py-3 font-mono text-xs text-gray-400 hover:text-gray-200">
                    {d.session_id.substring(0, 13)}...
                  </td>
                  <td className="px-4 py-3 text-gray-400 text-xs">
                    {formatDistanceToNow(new Date(d.timestamp), { addSuffix: true })}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          
          <div className="p-4 border-t border-gray-700 flex justify-between items-center text-sm">
            <span className="text-gray-400">Showing {data?.detections.length} results</span>
            <div className="space-x-2">
              <button className="px-3 py-1 bg-gray-700 hover:bg-gray-600 rounded text-gray-200 disabled:opacity-50">Prev</button>
              <button className="px-3 py-1 bg-gray-700 hover:bg-gray-600 rounded text-gray-200">Next</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

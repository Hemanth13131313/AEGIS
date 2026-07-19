import React, { useState } from 'react';
import { useRedTeamResults, useRedTeamSummary } from '../hooks/useRedTeamResults';

export default function RedTeamPage() {
  const { data: results, isLoading: isLoadingResults } = useRedTeamResults();
  const { data: summary, isLoading: isLoadingSummary } = useRedTeamSummary();
  
  const [owaspFilter, setOwaspFilter] = useState<string>('all');
  const [resultFilter, setResultFilter] = useState<string>('all');

  const filteredResults = results?.filter(r => {
    if (owaspFilter !== 'all' && r.owasp_category !== owaspFilter) return false;
    if (resultFilter === 'pass' && !r.passed) return false;
    if (resultFilter === 'fail' && r.passed) return false;
    if (resultFilter === 'disagreement' && !r.disagreement) return false;
    return true;
  });

  const getPassRateColor = (rate: number) => {
    if (rate > 0.9) return 'text-green-600';
    if (rate >= 0.7) return 'text-yellow-600';
    return 'text-red-600';
  };

  return (
    <div className="p-6">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">Red Team</h1>
          <p className="text-sm text-gray-500">Continuous adversarial validation of the detection pipeline</p>
        </div>
        <button 
          className="bg-blue-600 text-white px-4 py-2 rounded shadow hover:bg-blue-700 disabled:opacity-50"
          title="Configure AEGIS_REDTEAM_TARGET_URL to enable"
        >
          Run Now
        </button>
      </div>

      {!isLoadingSummary && summary && (
        <div className="grid grid-cols-4 gap-4 mb-6">
          <div className="bg-white p-4 rounded shadow border border-gray-100">
            <div className="text-sm text-gray-500">Test Cases</div>
            <div className="text-2xl font-bold">{summary.total}</div>
          </div>
          <div className="bg-white p-4 rounded shadow border border-gray-100">
            <div className="text-sm text-gray-500">Last Run</div>
            <div className="text-2xl font-bold">Just now</div>
          </div>
          <div className="bg-white p-4 rounded shadow border border-gray-100">
            <div className="text-sm text-gray-500">Pass Rate</div>
            <div className={`text-2xl font-bold ${getPassRateColor(summary.pass_rate)}`}>
              {(summary.pass_rate * 100).toFixed(0)}%
            </div>
          </div>
          <div className="bg-white p-4 rounded shadow border border-gray-100">
            <div className="text-sm text-gray-500">Disagreements</div>
            <div className="text-2xl font-bold text-orange-500">
              {summary.disagreements} <span className="text-xs bg-orange-100 text-orange-800 px-2 py-1 rounded-full">Alert</span>
            </div>
          </div>
        </div>
      )}

      <div className="bg-white rounded shadow border border-gray-200 overflow-hidden">
        <div className="p-4 border-b border-gray-200 flex gap-4 bg-gray-50">
          <select 
            value={owaspFilter} 
            onChange={(e) => setOwaspFilter(e.target.value)}
            className="border-gray-300 rounded shadow-sm"
          >
            <option value="all">All OWASP Categories</option>
            <option value="LLM01">LLM01 - Prompt Injection</option>
            <option value="LLM06">LLM06 - Data Exfiltration</option>
            <option value="none">Safe / Benign</option>
          </select>

          <select 
            value={resultFilter} 
            onChange={(e) => setResultFilter(e.target.value)}
            className="border-gray-300 rounded shadow-sm"
          >
            <option value="all">All Results</option>
            <option value="pass">Passed</option>
            <option value="fail">Failed</option>
            <option value="disagreement">Disagreements</option>
          </select>
        </div>

        <table className="min-w-full divide-y divide-gray-200 text-sm">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500 uppercase">Test ID</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500 uppercase">Name</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500 uppercase">Category</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500 uppercase">Technique</th>
              <th className="px-6 py-3 text-left font-medium text-gray-500 uppercase">Action</th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {!isLoadingResults && filteredResults?.map((r) => {
              let rowClass = "bg-white";
              let statusIcon = "❌";
              
              if (r.passed) {
                rowClass = r.disagreement ? "bg-orange-50" : "bg-green-50";
                statusIcon = r.disagreement ? "⚠️" : "✅";
              } else {
                rowClass = "bg-red-50";
              }

              return (
                <tr key={r.test_case_id} className={rowClass}>
                  <td className="px-6 py-4 whitespace-nowrap">{statusIcon}</td>
                  <td className="px-6 py-4 whitespace-nowrap font-mono">{r.test_case_id}</td>
                  <td className="px-6 py-4">{r.name}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    <span className="bg-gray-100 text-gray-800 px-2 py-1 rounded text-xs">{r.owasp_category}</span>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-xs text-gray-500">{r.atlas_technique}</td>
                  <td className="px-6 py-4 whitespace-nowrap">
                    Got: <b>{r.action_received}</b><br/>
                    <span className="text-xs text-gray-500">Exp: {r.expected_action}</span>
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
    </div>
  );
}

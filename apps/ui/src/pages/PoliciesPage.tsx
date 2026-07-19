import React, { useState } from 'react';
import { usePolicies, useCreatePolicy } from '../hooks/usePolicies';

export const PoliciesPage: React.FC = () => {
  const { data, isLoading } = usePolicies();
  const [isModalOpen, setModalOpen] = useState(false);
  const createPolicy = useCreatePolicy();

  return (
    <div className="p-6">
      <h1 className="text-2xl font-bold mb-4">Policies</h1>
      <div className="flex justify-between items-center mb-6">
        <div className="text-sm text-gray-500">
          Scope: Org → App → Model → Environment
        </div>
        <button 
          className="bg-blue-600 text-white px-4 py-2 rounded"
          onClick={() => setModalOpen(true)}
        >
          New Policy
        </button>
      </div>

      {isLoading ? (
        <div>Loading policies...</div>
      ) : data?.policies?.length === 0 ? (
        <div className="text-gray-500">No policies yet. Create your first policy to start enforcing AI traffic.</div>
      ) : (
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b">
              <th className="p-2">Scope Type</th>
              <th className="p-2">Status</th>
              <th className="p-2">Version</th>
              <th className="p-2">Last Modified</th>
              <th className="p-2">Actions</th>
            </tr>
          </thead>
          <tbody>
            {data?.policies.map((p) => (
              <tr key={p.id} className="border-b">
                <td className="p-2"><span className="px-2 py-1 bg-gray-100 rounded text-sm">{p.scope_type}</span></td>
                <td className="p-2 flex items-center gap-2">
                  <div className={`w-2 h-2 rounded-full ${p.active ? 'bg-green-500' : 'bg-red-500'}`}></div>
                  {p.active ? 'Active' : 'Inactive'}
                </td>
                <td className="p-2">v{p.version || 1}</td>
                <td className="p-2">{new Date(p.updated_at).toLocaleDateString()}</td>
                <td className="p-2 text-blue-600 cursor-pointer">View/Edit</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}

      {isModalOpen && (
        <div className="fixed inset-0 bg-black/50 flex items-center justify-center p-4">
          <div className="bg-white rounded-lg p-6 w-full max-w-2xl">
            <h2 className="text-xl font-bold mb-4">Create Policy</h2>
            <textarea 
              className="w-full h-48 border p-2 font-mono text-sm"
              placeholder="package aegis.policy..."
              defaultValue="package aegis.policy"
            ></textarea>
            <div className="flex justify-end gap-2 mt-4">
              <button onClick={() => setModalOpen(false)} className="px-4 py-2 border rounded">Cancel</button>
              <button className="px-4 py-2 bg-blue-600 text-white rounded">Submit</button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export interface Policy {
  id: string;
  scope_id: string;
  scope_type: 'organization' | 'application' | 'model_endpoint' | 'environment';
  rego_bundle_ref: string;
  active: boolean;
  created_at: string;
  updated_at: string;
  version?: number;
  rego_content?: string;
  actor?: string;
}

export interface CreatePolicyInput {
  scope_id: string;
  scope_type: string;
  rego_content: string;
  actor: string;
}

// TODO: Phase 4 replace with real API calls using react-query
export function usePolicies(cursor?: string) {
  return {
    data: {
      policies: [] as Policy[],
      next_cursor: ''
    },
    isLoading: false,
    error: null
  };
}

export function usePolicy(id: string) {
  return {
    data: null,
    isLoading: false,
    error: null
  };
}

export function useCreatePolicy() {
  return (input: CreatePolicyInput) => Promise.resolve({});
}

export function useUpdatePolicy() {
  return (id: string, content: string) => Promise.resolve({});
}

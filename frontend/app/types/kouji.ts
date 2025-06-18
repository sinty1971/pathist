// Extended KoujiEntry type with company and location names
export interface KoujiEntryExtended {
  id: number;
  name: string;
  path: string;
  is_directory: boolean;
  size: number;
  modified_time: string;
  company_name?: string;
  location_name?: string;
  status?: string;
  start_date?: string;
  end_date?: string;
  description?: string;
  tags?: string[];
}

export interface KoujiEntriesResponseExtended {
  entries: KoujiEntryExtended[];
  count: number;
  path: string;
  total_size?: number;
}
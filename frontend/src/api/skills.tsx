import axios from 'axios';

const API_URL = 'http://localhost:8080/api'; // URL вашего Go-бэкенда

export interface SkillStats {
  name: string;
  count: number;
}

export const fetchSkillsStats = async (limit: number = 25): Promise<SkillStats[]> => {
  const response = await axios.get<SkillStats[]>(`${API_URL}/skills/stats?limit=${limit}`);
  return response.data;
};
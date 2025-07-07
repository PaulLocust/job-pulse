import { useEffect, useState } from 'react';
import { fetchSkillsStats, type SkillStats } from '../api/skills';
import { SkillsChart } from '../components/SkillsChart';

export const Stats = () => {
  const [skills, setSkills] = useState<SkillStats[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadData = async () => {
      try {
        const data = await fetchSkillsStats(25);
        setSkills(data);
      } catch (err) {
        setError('Ошибка загрузки данных');
      } finally {
        setIsLoading(false);
      }
    };
    loadData();
  }, []);

  return (
    <div className="page">
      <h1 className="page-title">Топ навыков</h1>
      <div className="chart-container">
        {isLoading ? (
          <p className="loading-text">Загрузка данных...</p>
        ) : error ? (
          <p className="error-message">{error}</p>
        ) : (
          <div className="chart-wrapper">
            <SkillsChart data={skills} />
          </div>
        )}
      </div>
    </div>
  );
};
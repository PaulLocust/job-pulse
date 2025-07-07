import { useEffect, useState } from 'react';
import { fetchSkillsStats, type SkillStats } from '../api/skills';
import { SkillsChart } from '../components/SkillsChart';

export const Stats = () => {
  const [skills, setSkills] = useState<SkillStats[]>([]);

  useEffect(() => {
    const loadData = async () => {
      const data = await fetchSkillsStats(25);
      setSkills(data);
    };
    loadData();
  }, []);

  return (
    <div className="page">
      <h1>Топ навыков</h1>
      {skills.length > 0 ? (
        <div style={{ width: '800px', margin: '0 auto' }}>
          <SkillsChart data={skills} />
        </div>
      ) : (
        <p>Загрузка данных...</p>
      )}
    </div>
  );
};
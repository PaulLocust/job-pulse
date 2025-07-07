import { useEffect, useRef } from 'react';
import { Chart, registerables } from 'chart.js';

Chart.register(...registerables);

interface SkillsChartProps {
  data: { name: string; count: number }[];
}

export const SkillsChart = ({ data }: SkillsChartProps) => {
  const chartRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!chartRef.current || data.length === 0) return;

    const chart = new Chart(chartRef.current, {
      type: 'bar',
      data: {
        labels: data.map((skill) => skill.name),
        datasets: [{
          label: 'Упоминаний в вакансиях',
          data: data.map((skill) => skill.count),
          backgroundColor: 'rgba(54, 162, 235, 0.7)',
        }]
      },
      options: {
        indexAxis: 'y',
        responsive: true,
        maintainAspectRatio: false, // Отключаем авто-масштабирование
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            callbacks: {
              label: (ctx) => `${ctx.raw} упоминаний`,
            },
          },
        },
        scales: {
          y: {
            ticks: {
              mirror: true, // Размещаем подписи слева от столбцов
              padding: 10, // Отступ для названий
              font: {
                size: 12, // Размер шрифта
              },
            },
          },
          x: {
            beginAtZero: true,
            title: {
              display: true,
              text: 'Количество упоминаний',
            },
          },
        },
      },
    });

    return () => chart.destroy();
  }, [data]);

  return (
    <div style={{ 
      width: '100%', 
      height: `${data.length * 40}px`, // Динамическая высота
      minHeight: '400px',
    }}>
      <canvas ref={chartRef} />
    </div>
  );
};
import { useEffect, useRef } from 'react';
import { Chart, registerables } from 'chart.js';

// Регистрируем необходимые компоненты Chart.js
Chart.register(...registerables);

interface SkillsChartProps {
  data: {
    name: string;
    count: number;
  }[];
}

export const SkillsChart = ({ data }: SkillsChartProps) => {
  const chartRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    if (!chartRef.current || data.length === 0) return;

    const chart = new Chart(chartRef.current, {
      type: 'bar',
      data: {
        labels: data.map((skill) => skill.name),
        datasets: [
          {
            label: 'Количество упоминаний',
            data: data.map((skill) => skill.count),
            backgroundColor: 'rgba(54, 162, 235, 0.7)',
            borderColor: 'rgba(54, 162, 235, 1)',
            borderWidth: 1,
          },
        ],
      },
      options: {
        indexAxis: 'y',
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false,
          },
          tooltip: {
            callbacks: {
              label: (context) => `${context.raw} упоминаний`,
            },
          },
        },
        scales: {
          y: {
            ticks: {
              color: '#000000', // Черный цвет для названий навыков
              font: {
                size: 12,
                weight: 'bold',
              },
              padding: 10,
              mirror: true, // Отображаем текст слева от столбцов
            },
            grid: {
              display: false, // Скрываем сетку по оси Y
            },
          },
          x: {
            ticks: {
              color: '#000000', // Черный цвет для цифр
            },
            title: {
              display: true,
              text: 'Количество упоминаний',
              color: '#000000', // Черный цвет для заголовка
              font: {
                weight: 'bold',
              },
            },
            grid: {
              color: 'rgba(0, 0, 0, 0.1)', // Светло-серая сетка
            },
            beginAtZero: true,
          },
        },
      },
    });

    return () => chart.destroy();
  }, [data]);

  return (
    <div
      style={{
        width: '100%',
        height: `${Math.max(data.length * 40, 400)}px`, // Минимальная высота 400px
        position: 'relative',
      }}
    >
      <canvas ref={chartRef} />
    </div>
  );
};
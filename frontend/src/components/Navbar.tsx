import { Link } from 'react-router-dom';

export const Navbar = () => {
  return (
    <nav>
      <ul>
        <li>
          <Link to="/">Главная</Link>
        </li>
        <li>
          <Link to="/stats">Статистика</Link>
        </li>
      </ul>
    </nav>
  );
};
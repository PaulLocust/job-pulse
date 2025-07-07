import { Link } from 'react-router-dom';

export const Navbar = () => {
  return (
    <nav className="navbar">
      <ul className="nav-list">
        <li>
          <Link to="/" className="nav-link">Главная</Link>
        </li>
        <li>
          <Link to="/stats" className="nav-link">Статистика</Link>
        </li>
      </ul>
    </nav>
  );
};
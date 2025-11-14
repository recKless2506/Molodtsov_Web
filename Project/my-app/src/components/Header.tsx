// src/components/Header.tsx
import React from "react";
import CartIcon from "./CartIcon";
import "./header.css";

type HeaderProps = {
  cartCount?: number;
};

const Header: React.FC<HeaderProps> = ({ cartCount = 0 }) => {
  return (
    <header className="header">
      <div className="logo">
        <a href="/">
          {/* svg логотип */}
          <svg width="58" height="58" viewBox="0 0 58 58" fill="none" xmlns="http://www.w3.org/2000/svg">
            <g clipPath="url(#clip0)"><path d="M21.79 55.4667L21.7879 35.8916L34.2624 35.8903L34.2644 55.4654L49.8574 55.4638L49.8547 29.3637L59.2105 29.3627L28.0214 0.00324703L-3.16159 29.3692L6.19423 29.3682L6.19693 55.4683L21.79 55.4667Z" fill="white"/></g>
            <defs><clipPath id="clip0"><rect width="57.6862" height="57.6862" fill="white"/></clipPath></defs>
          </svg>
        </a>
      </div>

      <nav style={{ marginLeft: "auto", display: "flex", alignItems: "center", gap: 12 }}>
        <a href="/catalog" style={{ color: "white", textDecoration: "none", fontWeight: 600 }}>Каталог</a>
        <CartIcon count={cartCount} />
      </nav>
    </header>
  );
};

export default Header;

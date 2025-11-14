// src/pages/HomePage.tsx
import React from "react";
import Header from "../components/Header";
import "./home.css";

const HomePage: React.FC<{ cartCount: number }> = ({ cartCount }) => {
  return (
    <div>
      <Header cartCount={cartCount} />
      <div className="home-container">
        <h1>Добро пожаловать в магазин теплонагревателей</h1>
        <p>Выберите подходящий прибор в каталоге.</p>
        <a className="home-button" href="/catalog">Перейти в каталог</a>
      </div>
    </div>
  );
};

export default HomePage;

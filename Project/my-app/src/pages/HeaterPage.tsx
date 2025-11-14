// src/pages/HeaterPage.tsx
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import Header from "../components/Header";
import type { HeaterProduct } from "../types";
import { fetchHeaterById } from "../api/api";
import "./heater.css";
import defaultImage from "../assets/DefaultImage.jpg";

const HeaterPage: React.FC<{ products?: HeaterProduct[]; cartCount: number }> = ({ cartCount }) => {
  const { id } = useParams<{ id: string }>();
  const [product, setProduct] = useState<HeaterProduct | null>(null);

  useEffect(() => {
    if (!id) return;
    (async () => {
      try {
        // попытка получить продукт по API
        const data = await fetchHeaterById(id);
        // если бек отдаёт объект { product: {...} } или сам объект
        setProduct(data.product ?? data);
      } catch (e) {
        console.warn("Failed to load heater by id", e);
      }
    })();
  }, [id]);

  if (!product) return <div><Header cartCount={cartCount} /><div style={{ padding: 40, textAlign: "center" }}>Загрузка...</div></div>;

  return (
    <div>
      <Header cartCount={cartCount} />
      <div className="product-container">
        <div className="product-image">
          <img src={product.Image ?? defaultImage} alt={product.Title} />
        </div>
        <div className="product-content">
          <div className="product-title">{product.Title}</div>
          <div className="spec-title">Описание</div>
          <div className="product-description">{product.Description}</div>
          <div className="spec-title">Технические характеристики</div>
          <div className="spec-text">{product.Power}</div>
        </div>
      </div>
    </div>
  );
};

export default HeaterPage;

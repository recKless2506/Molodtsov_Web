// src/components/ProductCard.tsx
import React from "react";
import type { HeaterProduct } from "../types";
import "./productcard.css";
import defaultImage from "../assets/DefaultImage.jpg";

type Props = {
  product: HeaterProduct;
  onAddToCart: (p: HeaterProduct) => void;
};

const ProductCard: React.FC<Props> = ({ product, onAddToCart }) => {
  return (
    <div className="product-card">
      <img src={product.Image ?? defaultImage} alt={product.Title} />
      <div className="content">
        <div className="product-title">{product.Title}</div>
        <div className="product-specs">{product.Efficiency}</div>
        <div style={{ marginTop: "auto", display: "flex", gap: "10px", justifyContent: "center" }}>
          <a className="product-button" href={`/heater/${product.ID}`}>Подробнее</a>
          <button className="product-button" onClick={() => onAddToCart(product)}>+</button>
        </div>
      </div>
    </div>
  );
};

export default ProductCard;

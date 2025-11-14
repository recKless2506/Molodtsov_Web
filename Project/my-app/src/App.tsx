// src/App.tsx
import React, { useEffect, useState } from "react";
import { Routes, Route } from "react-router-dom";
import HomePage from "./pages/HomePage";
import CatalogPage from "./pages/CatalogPage";
import HeaterPage from "./pages/HeaterPage";
import ApplicationPage from "./pages/ApplicationPage";
import type { HeaterProduct, Request } from "./types";
import { fetchCatalog } from "./api/api";

const App: React.FC = () => {
  const [products, setProducts] = useState<HeaterProduct[]>([]);
  const [cartCount, setCartCount] = useState<number>(0);
  const [requests, setRequests] = useState<Request[]>([]);

  useEffect(() => {
    // загрузка каталога при старте
    (async () => {
      try {
        const data = await fetchCatalog();
        // backend может возвращать шаблон HTML — предполагаем JSON { products: [...] } или массив
        const prods = Array.isArray(data) ? data : data.products ?? [];
        setProducts(prods);
      } catch (e) {
        console.error("Failed to load catalog:", e);
      }
    })();
  }, []);

  const onAddToCart = (product: HeaterProduct) => {
    // простая локальная логика: увеличиваем счётчик
    setCartCount((s) => s + 1);
    // можно послать POST /add-to-cart/:id, но сервер ожидает form POST (html). Покажу пример:
    fetch(`${import.meta.env.VITE_API_URL ?? "http://localhost:8001"}/add-to-cart/${product.ID}`, {
      method: "POST",
    }).catch((err) => console.warn("add-to-cart failed", err));
  };

  const clearCart = () => {
    setCartCount(0);
    setRequests([]);
    fetch(`${import.meta.env.VITE_API_URL ?? "http://localhost:8001"}/clear-cart`, { method: "POST" })
      .catch((err) => console.warn("clear-cart failed", err));
  };

  return (
    <Routes>
      <Route path="/" element={<HomePage cartCount={cartCount} />} />
      <Route
        path="/catalog"
        element={<CatalogPage products={products} cartCount={cartCount} onAddToCart={onAddToCart} />}
      />
      <Route path="/heater/:id" element={<HeaterPage products={products} cartCount={cartCount} />} />
      <Route path="/heaters_application" element={<ApplicationPage requests={requests} clearCart={clearCart} />} />
    </Routes>
  );
};

export default App;

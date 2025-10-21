# Bitoracdn (Go Version)

**Bitoracdn Go** is the second version of the Bitoracdn project — a **Go-based CDN system** designed to explore distributed caching, concurrency, and performance at scale.

It’s a **progressive rewrite** of the [NodeJS version](https://github.com/kjxcodez/bitoracdn-node), using Go to implement the same architecture with system-level efficiency.

---

## ⚙️ Architecture Overview

````shell
User → Router → Edge Nodes → Origin Server
````

- **Origin** → File server / data source  
- **Edge** → Cache + serve content concurrently  
- **Router** → Load balances or geo-routes requests  
- **Control Plane** → Metrics, purge, health monitoring  

---

## Learning Goals

- Master **Go concurrency** with goroutines & channels  
- Build concurrent **edge caches** (LRU, TTL-based)  
- Implement **distributed invalidation** and **health checks**  
- Add **observability** using `pprof`, Prometheus, and logs  
- Compare **Node vs Go** performance in real-world scenarios  

---

## Tech Stack

- **Go 1.23+**
- **net/http**, **context**, **sync**, **prometheus/client_golang**
- **Docker**
- **GeoLite2** (for Geo-routing)

---

## Purpose

This project translates the concepts mastered in NodeJS into Go — focusing on **performance, type safety, concurrency, and system design**.

You’ll learn how distributed systems scale in a compiled, low-latency language while maintaining the same architecture.

---

## License

MIT License © 2025 [Kapil Jangid](https://github.com/kjxcodez)

## Socials

![X (formerly Twitter) Follow](https://img.shields.io/twitter/follow/kjxcodez?style=social)
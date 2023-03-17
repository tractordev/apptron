// lib/qtalk.min.js
var lr = Object.defineProperty;
var fr = (t, e) => {
  for (var r in e)
    lr(t, r, { get: e[r], enumerable: true });
};
var je;
try {
  je = new TextDecoder();
} catch {
}
var p;
var ie;
var l = 0;
var Rt = [];
var ur = 105;
var hr = 57342;
var yr = 57343;
var _t = 57337;
var Bt = 6;
var ce = {};
var Ke = Rt;
var He = 0;
var D = {};
var P;
var Ee;
var Se = 0;
var ye = 0;
var R;
var j;
var k = [];
var qe = [];
var T;
var W;
var pe;
var Ft = { useRecords: false, mapsAsObjects: true };
var we = false;
var Y = class {
  constructor(e) {
    if (e && ((e.keyMap || e._keyMap) && !e.useRecords && (e.useRecords = false, e.mapsAsObjects = true), e.useRecords === false && e.mapsAsObjects === void 0 && (e.mapsAsObjects = true), e.getStructures && (e.getShared = e.getStructures), e.getShared && !e.structures && ((e.structures = []).uninitialized = true), e.keyMap)) {
      this.mapKey = /* @__PURE__ */ new Map();
      for (let [r, n] of Object.entries(e.keyMap))
        this.mapKey.set(n, r);
    }
    Object.assign(this, e);
  }
  decodeKey(e) {
    return this.keyMap && this.mapKey.get(e) || e;
  }
  encodeKey(e) {
    return this.keyMap && this.keyMap.hasOwnProperty(e) ? this.keyMap[e] : e;
  }
  encodeKeys(e) {
    if (!this._keyMap)
      return e;
    let r = /* @__PURE__ */ new Map();
    for (let [n, s] of Object.entries(e))
      r.set(this._keyMap.hasOwnProperty(n) ? this._keyMap[n] : n, s);
    return r;
  }
  decodeKeys(e) {
    if (!this._keyMap || e.constructor.name != "Map")
      return e;
    if (!this._mapKey) {
      this._mapKey = /* @__PURE__ */ new Map();
      for (let [n, s] of Object.entries(this._keyMap))
        this._mapKey.set(s, n);
    }
    let r = {};
    return e.forEach((n, s) => r[v(this._mapKey.has(s) ? this._mapKey.get(s) : s)] = n), r;
  }
  mapDecode(e, r) {
    let n = this.decode(e);
    if (this._keyMap)
      switch (n.constructor.name) {
        case "Array":
          return n.map((s) => this.decodeKeys(s));
      }
    return n;
  }
  decode(e, r) {
    if (p)
      return zt(() => (Me(), this ? this.decode(e, r) : Y.prototype.decode.call(Ft, e, r)));
    ie = r > -1 ? r : e.length, l = 0, He = 0, ye = 0, Ee = null, Ke = Rt, R = null, p = e;
    try {
      W = e.dataView || (e.dataView = new DataView(e.buffer, e.byteOffset, e.byteLength));
    } catch (n) {
      throw p = null, e instanceof Uint8Array ? n : new Error("Source must be a Uint8Array or Buffer but was a " + (e && typeof e == "object" ? e.constructor.name : typeof e));
    }
    if (this instanceof Y) {
      if (D = this, T = this.sharedValues && (this.pack ? new Array(this.maxPrivatePackedValues || 16).concat(this.sharedValues) : this.sharedValues), this.structures)
        return P = this.structures, Ue();
      (!P || P.length > 0) && (P = []);
    } else
      D = Ft, (!P || P.length > 0) && (P = []), T = null;
    return Ue();
  }
  decodeMultiple(e, r) {
    let n, s = 0;
    try {
      let o = e.length;
      we = true;
      let f = this ? this.decode(e, o) : Xe.decode(e, o);
      if (r) {
        if (r(f) === false)
          return;
        for (; l < o; )
          if (s = l, r(Ue()) === false)
            return;
      } else {
        for (n = [f]; l < o; )
          s = l, n.push(Ue());
        return n;
      }
    } catch (o) {
      throw o.lastPosition = s, o.values = n, o;
    } finally {
      we = false, Me();
    }
  }
};
function Ue() {
  try {
    let t = E();
    if (R) {
      if (l >= R.postBundlePosition) {
        let e = new Error("Unexpected bundle position");
        throw e.incomplete = true, e;
      }
      l = R.postBundlePosition, R = null;
    }
    if (l == ie)
      P = null, p = null, j && (j = null);
    else if (l > ie) {
      let e = new Error("Unexpected end of CBOR data");
      throw e.incomplete = true, e;
    } else if (!we)
      throw new Error("Data read, but end of buffer not reached");
    return t;
  } catch (t) {
    throw Me(), (t instanceof RangeError || t.message.startsWith("Unexpected end of buffer")) && (t.incomplete = true), t;
  }
}
function E() {
  let t = p[l++], e = t >> 5;
  if (t = t & 31, t > 23)
    switch (t) {
      case 24:
        t = p[l++];
        break;
      case 25:
        if (e == 7)
          return xr();
        t = W.getUint16(l), l += 2;
        break;
      case 26:
        if (e == 7) {
          let r = W.getFloat32(l);
          if (D.useFloat32 > 2) {
            let n = Pe[(p[l] & 127) << 1 | p[l + 1] >> 7];
            return l += 4, (n * r + (r > 0 ? 0.5 : -0.5) >> 0) / n;
          }
          return l += 4, r;
        }
        t = W.getUint32(l), l += 4;
        break;
      case 27:
        if (e == 7) {
          let r = W.getFloat64(l);
          return l += 8, r;
        }
        if (e > 1) {
          if (W.getUint32(l) > 0)
            throw new Error("JavaScript does not support arrays, maps, or strings with length over 4294967295");
          t = W.getUint32(l + 4);
        } else
          D.int64AsNumber ? (t = W.getUint32(l) * 4294967296, t += W.getUint32(l + 4)) : t = W.getBigUint64(l);
        l += 8;
        break;
      case 31:
        switch (e) {
          case 2:
          case 3:
            throw new Error("Indefinite length not supported for byte or text strings");
          case 4:
            let r = [], n, s = 0;
            for (; (n = E()) != ce; )
              r[s++] = n;
            return e == 4 ? r : e == 3 ? r.join("") : Buffer.concat(r);
          case 5:
            let o;
            if (D.mapsAsObjects) {
              let f = {};
              if (D.keyMap)
                for (; (o = E()) != ce; )
                  f[v(D.decodeKey(o))] = E();
              else
                for (; (o = E()) != ce; )
                  f[v(o)] = E();
              return f;
            } else {
              pe && (D.mapsAsObjects = true, pe = false);
              let f = /* @__PURE__ */ new Map();
              if (D.keyMap)
                for (; (o = E()) != ce; )
                  f.set(D.decodeKey(o), E());
              else
                for (; (o = E()) != ce; )
                  f.set(o, E());
              return f;
            }
          case 7:
            return ce;
          default:
            throw new Error("Invalid major type for indefinite length " + e);
        }
      default:
        throw new Error("Unknown token " + t);
    }
  switch (e) {
    case 0:
      return t;
    case 1:
      return ~t;
    case 2:
      return mr(t);
    case 3:
      if (ye >= l)
        return Ee.slice(l - Se, (l += t) - Se);
      if (ye == 0 && ie < 140 && t < 32) {
        let s = t < 16 ? Tt(t) : wr(t);
        if (s != null)
          return s;
      }
      return pr(t);
    case 4:
      let r = new Array(t);
      for (let s = 0; s < t; s++)
        r[s] = E();
      return r;
    case 5:
      if (D.mapsAsObjects) {
        let s = {};
        if (D.keyMap)
          for (let o = 0; o < t; o++)
            s[v(D.decodeKey(E()))] = E();
        else
          for (let o = 0; o < t; o++)
            s[v(E())] = E();
        return s;
      } else {
        pe && (D.mapsAsObjects = true, pe = false);
        let s = /* @__PURE__ */ new Map();
        if (D.keyMap)
          for (let o = 0; o < t; o++)
            s.set(D.decodeKey(E()), E());
        else
          for (let o = 0; o < t; o++)
            s.set(E(), E());
        return s;
      }
    case 6:
      if (t >= _t) {
        let s = P[t & 8191];
        if (s)
          return s.read || (s.read = $e(s)), s.read();
        if (t < 65536) {
          if (t == yr) {
            let o = fe(), f = E(), m = E();
            Je(f, m);
            let g = {};
            if (D.keyMap)
              for (let w = 2; w < o; w++) {
                let O = D.decodeKey(m[w - 2]);
                g[v(O)] = E();
              }
            else
              for (let w = 2; w < o; w++) {
                let O = m[w - 2];
                g[v(O)] = E();
              }
            return g;
          } else if (t == hr) {
            let o = fe(), f = E();
            for (let m = 2; m < o; m++)
              Je(f++, E());
            return E();
          } else if (t == _t)
            return Er();
          if (D.getShared && (Ye(), s = P[t & 8191], s))
            return s.read || (s.read = $e(s)), s.read();
        }
      }
      let n = k[t];
      if (n)
        return n.handlesRead ? n(E) : n(E());
      {
        let s = E();
        for (let o = 0; o < qe.length; o++) {
          let f = qe[o](t, s);
          if (f !== void 0)
            return f;
        }
        return new K(s, t);
      }
    case 7:
      switch (t) {
        case 20:
          return false;
        case 21:
          return true;
        case 22:
          return null;
        case 23:
          return;
        case 31:
        default:
          let s = (T || oe())[t];
          if (s !== void 0)
            return s;
          throw new Error("Unknown token " + t);
      }
    default:
      if (isNaN(t)) {
        let s = new Error("Unexpected end of CBOR data");
        throw s.incomplete = true, s;
      }
      throw new Error("Unknown CBOR token " + t);
  }
}
var Wt = /^[a-zA-Z_$][a-zA-Z\d_$]*$/;
function $e(t) {
  function e() {
    let r = p[l++];
    if (r = r & 31, r > 23)
      switch (r) {
        case 24:
          r = p[l++];
          break;
        case 25:
          r = W.getUint16(l), l += 2;
          break;
        case 26:
          r = W.getUint32(l), l += 4;
          break;
        default:
          throw new Error("Expected array header, but got " + p[l - 1]);
      }
    let n = this.compiledReader;
    for (; n; ) {
      if (n.propertyCount === r)
        return n(E);
      n = n.next;
    }
    if (this.slowReads++ >= 3) {
      let o = this.length == r ? this : this.slice(0, r);
      return n = D.keyMap ? new Function("r", "return {" + o.map((f) => D.decodeKey(f)).map((f) => Wt.test(f) ? v(f) + ":r()" : "[" + JSON.stringify(f) + "]:r()").join(",") + "}") : new Function("r", "return {" + o.map((f) => Wt.test(f) ? v(f) + ":r()" : "[" + JSON.stringify(f) + "]:r()").join(",") + "}"), this.compiledReader && (n.next = this.compiledReader), n.propertyCount = r, this.compiledReader = n, n(E);
    }
    let s = {};
    if (D.keyMap)
      for (let o = 0; o < r; o++)
        s[v(D.decodeKey(this[o]))] = E();
    else
      for (let o = 0; o < r; o++)
        s[v(this[o])] = E();
    return s;
  }
  return t.slowReads = 0, e;
}
function v(t) {
  return t === "__proto__" ? "__proto_" : t;
}
var pr = Ge;
function Ge(t) {
  let e;
  if (t < 16 && (e = Tt(t)))
    return e;
  if (t > 64 && je)
    return je.decode(p.subarray(l, l += t));
  let r = l + t, n = [];
  for (e = ""; l < r; ) {
    let s = p[l++];
    if ((s & 128) === 0)
      n.push(s);
    else if ((s & 224) === 192) {
      let o = p[l++] & 63;
      n.push((s & 31) << 6 | o);
    } else if ((s & 240) === 224) {
      let o = p[l++] & 63, f = p[l++] & 63;
      n.push((s & 31) << 12 | o << 6 | f);
    } else if ((s & 248) === 240) {
      let o = p[l++] & 63, f = p[l++] & 63, m = p[l++] & 63, g = (s & 7) << 18 | o << 12 | f << 6 | m;
      g > 65535 && (g -= 65536, n.push(g >>> 10 & 1023 | 55296), g = 56320 | g & 1023), n.push(g);
    } else
      n.push(s);
    n.length >= 4096 && (e += B.apply(String, n), n.length = 0);
  }
  return n.length > 0 && (e += B.apply(String, n)), e;
}
var B = String.fromCharCode;
function wr(t) {
  let e = l, r = new Array(t);
  for (let n = 0; n < t; n++) {
    let s = p[l++];
    if ((s & 128) > 0) {
      l = e;
      return;
    }
    r[n] = s;
  }
  return B.apply(String, r);
}
function Tt(t) {
  if (t < 4)
    if (t < 2) {
      if (t === 0)
        return "";
      {
        let e = p[l++];
        if ((e & 128) > 1) {
          l -= 1;
          return;
        }
        return B(e);
      }
    } else {
      let e = p[l++], r = p[l++];
      if ((e & 128) > 0 || (r & 128) > 0) {
        l -= 2;
        return;
      }
      if (t < 3)
        return B(e, r);
      let n = p[l++];
      if ((n & 128) > 0) {
        l -= 3;
        return;
      }
      return B(e, r, n);
    }
  else {
    let e = p[l++], r = p[l++], n = p[l++], s = p[l++];
    if ((e & 128) > 0 || (r & 128) > 0 || (n & 128) > 0 || (s & 128) > 0) {
      l -= 4;
      return;
    }
    if (t < 6) {
      if (t === 4)
        return B(e, r, n, s);
      {
        let o = p[l++];
        if ((o & 128) > 0) {
          l -= 5;
          return;
        }
        return B(e, r, n, s, o);
      }
    } else if (t < 8) {
      let o = p[l++], f = p[l++];
      if ((o & 128) > 0 || (f & 128) > 0) {
        l -= 6;
        return;
      }
      if (t < 7)
        return B(e, r, n, s, o, f);
      let m = p[l++];
      if ((m & 128) > 0) {
        l -= 7;
        return;
      }
      return B(e, r, n, s, o, f, m);
    } else {
      let o = p[l++], f = p[l++], m = p[l++], g = p[l++];
      if ((o & 128) > 0 || (f & 128) > 0 || (m & 128) > 0 || (g & 128) > 0) {
        l -= 8;
        return;
      }
      if (t < 10) {
        if (t === 8)
          return B(e, r, n, s, o, f, m, g);
        {
          let w = p[l++];
          if ((w & 128) > 0) {
            l -= 9;
            return;
          }
          return B(e, r, n, s, o, f, m, g, w);
        }
      } else if (t < 12) {
        let w = p[l++], O = p[l++];
        if ((w & 128) > 0 || (O & 128) > 0) {
          l -= 10;
          return;
        }
        if (t < 11)
          return B(e, r, n, s, o, f, m, g, w, O);
        let C = p[l++];
        if ((C & 128) > 0) {
          l -= 11;
          return;
        }
        return B(e, r, n, s, o, f, m, g, w, O, C);
      } else {
        let w = p[l++], O = p[l++], C = p[l++], L = p[l++];
        if ((w & 128) > 0 || (O & 128) > 0 || (C & 128) > 0 || (L & 128) > 0) {
          l -= 12;
          return;
        }
        if (t < 14) {
          if (t === 12)
            return B(e, r, n, s, o, f, m, g, w, O, C, L);
          {
            let F = p[l++];
            if ((F & 128) > 0) {
              l -= 13;
              return;
            }
            return B(e, r, n, s, o, f, m, g, w, O, C, L, F);
          }
        } else {
          let F = p[l++], z = p[l++];
          if ((F & 128) > 0 || (z & 128) > 0) {
            l -= 14;
            return;
          }
          if (t < 15)
            return B(e, r, n, s, o, f, m, g, w, O, C, L, F, z);
          let G = p[l++];
          if ((G & 128) > 0) {
            l -= 15;
            return;
          }
          return B(e, r, n, s, o, f, m, g, w, O, C, L, F, z, G);
        }
      }
    }
  }
}
function mr(t) {
  return D.copyBuffers ? Uint8Array.prototype.slice.call(p, l, l += t) : p.subarray(l, l += t);
}
var Vt = new Float32Array(1);
var Ce = new Uint8Array(Vt.buffer, 0, 4);
function xr() {
  let t = p[l++], e = p[l++], r = (t & 127) >> 2;
  if (r === 31)
    return e || t & 3 ? NaN : t & 128 ? -1 / 0 : 1 / 0;
  if (r === 0) {
    let n = ((t & 3) << 8 | e) / (1 << 24);
    return t & 128 ? -n : n;
  }
  return Ce[3] = t & 128 | (r >> 1) + 56, Ce[2] = (t & 7) << 5 | e >> 3, Ce[1] = e << 5, Ce[0] = 0, Vt[0];
}
var qr = new Array(4096);
var K = class {
  constructor(e, r) {
    this.value = e, this.tag = r;
  }
};
k[0] = (t) => new Date(t);
k[1] = (t) => new Date(Math.round(t * 1e3));
k[2] = (t) => {
  let e = BigInt(0);
  for (let r = 0, n = t.byteLength; r < n; r++)
    e = BigInt(t[r]) + e << BigInt(8);
  return e;
};
k[3] = (t) => BigInt(-1) - k[2](t);
k[4] = (t) => +(t[1] + "e" + t[0]);
k[5] = (t) => t[1] * Math.exp(t[0] * Math.log(2));
var Je = (t, e) => {
  t = t - 57344;
  let r = P[t];
  r && r.isShared && ((P.restoreStructures || (P.restoreStructures = []))[t] = r), P[t] = e, e.read = $e(e);
};
k[ur] = (t) => {
  let e = t.length, r = t[1];
  Je(t[0], r);
  let n = {};
  for (let s = 2; s < e; s++) {
    let o = r[s - 2];
    n[v(o)] = t[s];
  }
  return n;
};
k[14] = (t) => R ? R[0].slice(R.position0, R.position0 += t) : new K(t, 14);
k[15] = (t) => R ? R[1].slice(R.position1, R.position1 += t) : new K(t, 15);
var gr = { Error, RegExp };
k[27] = (t) => (gr[t[0]] || Error)(t[1], t[2]);
var Lt = (t) => {
  if (p[l++] != 132)
    throw new Error("Packed values structure must be followed by a 4 element array");
  let e = t();
  return T = T ? e.concat(T.slice(e.length)) : e, T.prefixes = t(), T.suffixes = t(), t();
};
Lt.handlesRead = true;
k[51] = Lt;
k[Bt] = (t) => {
  if (!T)
    if (D.getShared)
      Ye();
    else
      return new K(t, Bt);
  if (typeof t == "number")
    return T[16 + (t >= 0 ? 2 * t : -2 * t - 1)];
  throw new Error("No support for non-integer packed references yet");
};
k[28] = (t) => {
  j || (j = /* @__PURE__ */ new Map(), j.id = 0);
  let e = j.id++, r = p[l], n;
  r >> 5 == 4 ? n = [] : n = {};
  let s = { target: n };
  j.set(e, s);
  let o = t();
  return s.used ? Object.assign(n, o) : (s.target = o, o);
};
k[28].handlesRead = true;
k[29] = (t) => {
  let e = j.get(t);
  return e.used = true, e.target;
};
k[258] = (t) => new Set(t);
(k[259] = (t) => (D.mapsAsObjects && (D.mapsAsObjects = false, pe = true), t())).handlesRead = true;
function le(t, e) {
  return typeof t == "string" ? t + e : t instanceof Array ? t.concat(e) : Object.assign({}, t, e);
}
function oe() {
  if (!T)
    if (D.getShared)
      Ye();
    else
      throw new Error("No packed values available");
  return T;
}
var br = 1399353956;
qe.push((t, e) => {
  if (t >= 225 && t <= 255)
    return le(oe().prefixes[t - 224], e);
  if (t >= 28704 && t <= 32767)
    return le(oe().prefixes[t - 28672], e);
  if (t >= 1879052288 && t <= 2147483647)
    return le(oe().prefixes[t - 1879048192], e);
  if (t >= 216 && t <= 223)
    return le(e, oe().suffixes[t - 216]);
  if (t >= 27647 && t <= 28671)
    return le(e, oe().suffixes[t - 27639]);
  if (t >= 1811940352 && t <= 1879048191)
    return le(e, oe().suffixes[t - 1811939328]);
  if (t == br)
    return { packedValues: T, structures: P.slice(0), version: e };
  if (t == 55799)
    return e;
});
var Ir = new Uint8Array(new Uint16Array([1]).buffer)[0] == 1;
var vt = [Uint8Array, Uint8ClampedArray, Uint16Array, Uint32Array, typeof BigUint64Array == "undefined" ? { name: "BigUint64Array" } : BigUint64Array, Int8Array, Int16Array, Int32Array, typeof BigInt64Array == "undefined" ? { name: "BigInt64Array" } : BigInt64Array, Float32Array, Float64Array];
var Ar = [64, 68, 69, 70, 71, 72, 77, 78, 79, 85, 86];
for (let t = 0; t < vt.length; t++)
  Dr(vt[t], Ar[t]);
function Dr(t, e) {
  let r = "get" + t.name.slice(0, -5);
  typeof t != "function" && (t = null);
  let n = t.BYTES_PER_ELEMENT;
  for (let s = 0; s < 2; s++) {
    if (!s && n == 1)
      continue;
    let o = n == 2 ? 1 : n == 4 ? 2 : 3;
    k[s ? e : e - 4] = n == 1 || s == Ir ? (f) => {
      if (!t)
        throw new Error("Could not find typed array for code " + e);
      return new t(Uint8Array.prototype.slice.call(f, 0).buffer);
    } : (f) => {
      if (!t)
        throw new Error("Could not find typed array for code " + e);
      let m = new DataView(f.buffer, f.byteOffset, f.byteLength), g = f.length >> o, w = new t(g), O = m[r];
      for (let C = 0; C < g; C++)
        w[C] = O.call(m, C << o, s);
      return w;
    };
  }
}
function Er() {
  let t = fe(), e = l + E();
  for (let n = 2; n < t; n++)
    l += fe();
  let r = l;
  return l = e, R = [Ge(fe()), Ge(fe())], R.position0 = 0, R.position1 = 0, R.postBundlePosition = l, l = r, E();
}
function fe() {
  let t = p[l++] & 31;
  if (t > 23)
    switch (t) {
      case 24:
        t = p[l++];
        break;
      case 25:
        t = W.getUint16(l), l += 2;
        break;
      case 26:
        t = W.getUint32(l), l += 4;
        break;
    }
  return t;
}
function Ye() {
  if (D.getShared) {
    let t = zt(() => (p = null, D.getShared())) || {}, e = t.structures || [];
    D.sharedVersion = t.version, T = D.sharedValues = t.packedValues, P === true ? D.structures = P = e : P.splice.apply(P, [0, e.length].concat(e));
  }
}
function zt(t) {
  let e = ie, r = l, n = He, s = Se, o = ye, f = Ee, m = Ke, g = j, w = R, O = new Uint8Array(p.slice(0, ie)), C = P, L = D, F = we, z = t();
  return ie = e, l = r, He = n, Se = s, ye = o, Ee = f, Ke = m, j = g, R = w, p = O, we = F, P = C, D = L, W = new DataView(p.buffer, p.byteOffset, p.byteLength), z;
}
function Me() {
  p = null, j = null, P = null;
}
var Pe = new Array(147);
for (let t = 0; t < 256; t++)
  Pe[t] = +("1e" + Math.floor(45.15 - t * 0.30103));
var Xe = new Y({ useRecords: false });
var Ze = Xe.decode;
var Sr = Xe.decodeMultiple;
var Oe;
try {
  Oe = new TextEncoder();
} catch {
}
var Qe;
var jt;
var ke = globalThis.Buffer;
var me = typeof ke != "undefined";
var et = me ? ke.allocUnsafeSlow : Uint8Array;
var Kt = me ? ke : Uint8Array;
var Ht = 256;
var qt = me ? 4294967296 : 2144337920;
var tt;
var a;
var M;
var i = 0;
var te;
var _ = null;
var Ur = 61440;
var Cr = /[\u0080-\uFFFF]/;
var V = Symbol("record-id");
var Re = class extends Y {
  constructor(e) {
    super(e);
    this.offset = 0;
    let r, n, s, o, f, m;
    e = e || {};
    let g = Kt.prototype.utf8Write ? function(c, y, d) {
      return a.utf8Write(c, y, d);
    } : Oe && Oe.encodeInto ? function(c, y) {
      return Oe.encodeInto(c, a.subarray(y)).written;
    } : false, w = this, O = e.structures || e.saveStructures, C = e.maxSharedStructures;
    if (C == null && (C = O ? 128 : 0), C > 8190)
      throw new Error("Maximum maxSharedStructure is 8190");
    let L = e.sequential;
    L && (C = 0), this.structures || (this.structures = []), this.saveStructures && (this.saveShared = this.saveStructures);
    let F, z, G = e.sharedValues, N;
    if (G) {
      N = /* @__PURE__ */ Object.create(null);
      for (let c = 0, y = G.length; c < y; c++)
        N[G[c]] = c;
    }
    let J = [], ve = 0, De = 0;
    this.mapEncode = function(c, y) {
      if (this._keyMap && !this._mapped)
        switch (c.constructor.name) {
          case "Array":
            c = c.map((d) => this.encodeKeys(d));
            break;
        }
      return this.encode(c, y);
    }, this.encode = function(c, y) {
      if (a || (a = new et(8192), M = new DataView(a.buffer, 0, 8192), i = 0), te = a.length - 10, te - i < 2048 ? (a = new et(a.length), M = new DataView(a.buffer, 0, a.length), te = a.length - 10, i = 0) : y === ot && (i = i + 7 & 2147483640), n = i, w.useSelfDescribedHeader && (M.setUint32(i, 3654940416), i += 3), m = w.structuredClone ? /* @__PURE__ */ new Map() : null, w.bundleStrings && typeof c != "string" ? (_ = [], _.size = 1 / 0) : _ = null, s = w.structures, s) {
        if (s.uninitialized) {
          let h = w.getShared() || {};
          w.structures = s = h.structures || [], w.sharedVersion = h.version;
          let u = w.sharedValues = h.packedValues;
          if (u) {
            N = {};
            for (let x = 0, b = u.length; x < b; x++)
              N[u[x]] = x;
          }
        }
        let d = s.length;
        if (d > C && !L && (d = C), !s.transitions) {
          s.transitions = /* @__PURE__ */ Object.create(null);
          for (let h = 0; h < d; h++) {
            let u = s[h];
            if (!u)
              continue;
            let x, b = s.transitions;
            for (let I = 0, A = u.length; I < A; I++) {
              b[V] === void 0 && (b[V] = h);
              let S = u[I];
              x = b[S], x || (x = b[S] = /* @__PURE__ */ Object.create(null)), b = x;
            }
            b[V] = h | 1048576;
          }
        }
        L || (s.nextId = d);
      }
      if (o && (o = false), f = s || [], z = N, e.pack) {
        let d = /* @__PURE__ */ new Map();
        if (d.values = [], d.encoder = w, d.maxValues = e.maxPrivatePackedValues || (N ? 16 : 1 / 0), d.objectMap = N || false, d.samplingPackedValues = F, _e(c, d), d.values.length > 0) {
          a[i++] = 216, a[i++] = 51, H(4);
          let h = d.values;
          U(h), H(0), H(0), z = Object.create(N || null);
          for (let u = 0, x = h.length; u < x; u++)
            z[h[u]] = u;
        }
      }
      tt = y & at;
      try {
        if (tt)
          return;
        if (U(c), _ && Jt(n, U), w.offset = i, m && m.idsToInsert) {
          i += m.idsToInsert.length * 2, i > te && ue(i), w.offset = i;
          let d = Or(a.subarray(n, i), m.idsToInsert);
          return m = null, d;
        }
        return y & ot ? (a.start = n, a.end = i, a) : a.subarray(n, i);
      } finally {
        if (s) {
          if (De < 10 && De++, s.length > C && (s.length = C), ve > 1e4)
            s.transitions = null, De = 0, ve = 0, J.length > 0 && (J = []);
          else if (J.length > 0 && !L) {
            for (let d = 0, h = J.length; d < h; d++)
              J[d][V] = void 0;
            J = [];
          }
        }
        if (o && w.saveShared) {
          w.structures.length > C && (w.structures = w.structures.slice(0, C));
          let d = a.subarray(n, i);
          return w.updateSharedData() === false ? w.encode(c) : d;
        }
        y & _r && (i = n);
      }
    }, this.findCommonStringsToPack = () => (F = /* @__PURE__ */ new Map(), N || (N = /* @__PURE__ */ Object.create(null)), (c) => {
      let y = c && c.threshold || 4, d = this.pack ? c.maxPrivatePackedValues || 16 : 0;
      G || (G = this.sharedValues = []);
      for (let [h, u] of F)
        u.count > y && (N[h] = d++, G.push(h), o = true);
      for (; this.saveShared && this.updateSharedData() === false; )
        ;
      F = null;
    });
    let U = (c) => {
      i > te && (a = ue(i));
      var y = typeof c, d;
      if (y === "string") {
        if (z) {
          let b = z[c];
          if (b >= 0) {
            b < 16 ? a[i++] = b + 224 : (a[i++] = 198, b & 1 ? U(15 - b >> 1) : U(b - 16 >> 1));
            return;
          } else if (F && !e.pack) {
            let I = F.get(c);
            I ? I.count++ : F.set(c, { count: 1 });
          }
        }
        let h = c.length;
        if (_ && h >= 4 && h < 1024) {
          if ((_.size += h) > Ur) {
            let I, A = (_[0] ? _[0].length * 3 + _[1].length : 0) + 10;
            i + A > te && (a = ue(i + A)), a[i++] = 217, a[i++] = 223, a[i++] = 249, a[i++] = _.position ? 132 : 130, a[i++] = 26, I = i - n, i += 4, _.position && Jt(n, U), _ = ["", ""], _.size = 0, _.position = I;
          }
          let b = Cr.test(c);
          _[b ? 0 : 1] += c, a[i++] = b ? 206 : 207, U(h);
          return;
        }
        let u;
        h < 32 ? u = 1 : h < 256 ? u = 2 : h < 65536 ? u = 3 : u = 5;
        let x = h * 3;
        if (i + x > te && (a = ue(i + x)), h < 64 || !g) {
          let b, I, A, S = i + u;
          for (b = 0; b < h; b++)
            I = c.charCodeAt(b), I < 128 ? a[S++] = I : I < 2048 ? (a[S++] = I >> 6 | 192, a[S++] = I & 63 | 128) : (I & 64512) === 55296 && ((A = c.charCodeAt(b + 1)) & 64512) === 56320 ? (I = 65536 + ((I & 1023) << 10) + (A & 1023), b++, a[S++] = I >> 18 | 240, a[S++] = I >> 12 & 63 | 128, a[S++] = I >> 6 & 63 | 128, a[S++] = I & 63 | 128) : (a[S++] = I >> 12 | 224, a[S++] = I >> 6 & 63 | 128, a[S++] = I & 63 | 128);
          d = S - i - u;
        } else
          d = g(c, i + u, x);
        d < 24 ? a[i++] = 96 | d : d < 256 ? (u < 2 && a.copyWithin(i + 2, i + 1, i + 1 + d), a[i++] = 120, a[i++] = d) : d < 65536 ? (u < 3 && a.copyWithin(i + 3, i + 2, i + 2 + d), a[i++] = 121, a[i++] = d >> 8, a[i++] = d & 255) : (u < 5 && a.copyWithin(i + 5, i + 3, i + 3 + d), a[i++] = 122, M.setUint32(i, d), i += 4), i += d;
      } else if (y === "number")
        if (!this.alwaysUseFloat && c >>> 0 === c)
          c < 24 ? a[i++] = c : c < 256 ? (a[i++] = 24, a[i++] = c) : c < 65536 ? (a[i++] = 25, a[i++] = c >> 8, a[i++] = c & 255) : (a[i++] = 26, M.setUint32(i, c), i += 4);
        else if (!this.alwaysUseFloat && c >> 0 === c)
          c >= -24 ? a[i++] = 31 - c : c >= -256 ? (a[i++] = 56, a[i++] = ~c) : c >= -65536 ? (a[i++] = 57, M.setUint16(i, ~c), i += 2) : (a[i++] = 58, M.setUint32(i, ~c), i += 4);
        else {
          let h;
          if ((h = this.useFloat32) > 0 && c < 4294967296 && c >= -2147483648) {
            a[i++] = 250, M.setFloat32(i, c);
            let u;
            if (h < 4 || (u = c * Pe[(a[i] & 127) << 1 | a[i + 1] >> 7]) >> 0 === u) {
              i += 4;
              return;
            } else
              i--;
          }
          a[i++] = 251, M.setFloat64(i, c), i += 8;
        }
      else if (y === "object")
        if (!c)
          a[i++] = 246;
        else {
          if (m) {
            let u = m.get(c);
            if (u) {
              if (a[i++] = 216, a[i++] = 29, a[i++] = 25, !u.references) {
                let x = m.idsToInsert || (m.idsToInsert = []);
                u.references = [], x.push(u);
              }
              u.references.push(i - n), i += 2;
              return;
            } else
              m.set(c, { offset: i - n });
          }
          let h = c.constructor;
          if (h === Object)
            ze(c, true);
          else if (h === Array) {
            d = c.length, d < 24 ? a[i++] = 128 | d : H(d);
            for (let u = 0; u < d; u++)
              U(c[u]);
          } else if (h === Map)
            if ((this.mapsAsObjects ? this.useTag259ForMaps !== false : this.useTag259ForMaps) && (a[i++] = 217, a[i++] = 1, a[i++] = 3), d = c.size, d < 24 ? a[i++] = 160 | d : d < 256 ? (a[i++] = 184, a[i++] = d) : d < 65536 ? (a[i++] = 185, a[i++] = d >> 8, a[i++] = d & 255) : (a[i++] = 186, M.setUint32(i, d), i += 4), w.keyMap)
              for (let [u, x] of c)
                U(w.encodeKey(u)), U(x);
            else
              for (let [u, x] of c)
                U(u), U(x);
          else {
            for (let u = 0, x = Qe.length; u < x; u++) {
              let b = jt[u];
              if (c instanceof b) {
                let I = Qe[u], A = I.tag;
                A == null && (A = I.getTag && I.getTag.call(this, c)), A < 24 ? a[i++] = 192 | A : A < 256 ? (a[i++] = 216, a[i++] = A) : A < 65536 ? (a[i++] = 217, a[i++] = A >> 8, a[i++] = A & 255) : A > -1 && (a[i++] = 218, M.setUint32(i, A), i += 4), I.encode.call(this, c, U, ue);
                return;
              }
            }
            if (c[Symbol.iterator]) {
              if (tt) {
                let u = new Error("Iterable should be serialized as iterator");
                throw u.iteratorNotHandled = true, u;
              }
              a[i++] = 159;
              for (let u of c)
                U(u);
              a[i++] = 255;
              return;
            }
            if (c[Symbol.asyncIterator] || nt(c)) {
              let u = new Error("Iterable/blob should be serialized as iterator");
              throw u.iteratorNotHandled = true, u;
            }
            ze(c, !c.hasOwnProperty);
          }
        }
      else if (y === "boolean")
        a[i++] = c ? 245 : 244;
      else if (y === "bigint") {
        if (c < BigInt(1) << BigInt(64) && c >= 0)
          a[i++] = 27, M.setBigUint64(i, c);
        else if (c > -(BigInt(1) << BigInt(64)) && c < 0)
          a[i++] = 59, M.setBigUint64(i, -c - BigInt(1));
        else if (this.largeBigIntToFloat)
          a[i++] = 251, M.setFloat64(i, Number(c));
        else
          throw new RangeError(c + " was too large to fit in CBOR 64-bit integer format, set largeBigIntToFloat to convert to float-64");
        i += 8;
      } else if (y === "undefined")
        a[i++] = 247;
      else
        throw new Error("Unknown type: " + y);
    }, ze = this.useRecords === false ? this.variableMapSize ? (c) => {
      let y = Object.keys(c), d = Object.values(c), h = y.length;
      h < 24 ? a[i++] = 160 | h : h < 256 ? (a[i++] = 184, a[i++] = h) : h < 65536 ? (a[i++] = 185, a[i++] = h >> 8, a[i++] = h & 255) : (a[i++] = 186, M.setUint32(i, h), i += 4);
      let u;
      if (w.keyMap)
        for (let x = 0; x < h; x++)
          U(encodeKey(y[x])), U(d[x]);
      else
        for (let x = 0; x < h; x++)
          U(y[x]), U(d[x]);
    } : (c, y) => {
      a[i++] = 185;
      let d = i - n;
      i += 2;
      let h = 0;
      if (w.keyMap)
        for (let u in c)
          (y || c.hasOwnProperty(u)) && (U(w.encodeKey(u)), U(c[u]), h++);
      else
        for (let u in c)
          (y || c.hasOwnProperty(u)) && (U(u), U(c[u]), h++);
      a[d++ + n] = h >> 8, a[d + n] = h & 255;
    } : (c, y) => {
      let d, h = f.transitions || (f.transitions = /* @__PURE__ */ Object.create(null)), u = 0, x = 0, b, I;
      if (this.keyMap) {
        I = Object.keys(c).map((S) => this.encodeKey(S)), x = I.length;
        for (let S = 0; S < x; S++) {
          let Pt = I[S];
          d = h[Pt], d || (d = h[Pt] = /* @__PURE__ */ Object.create(null), u++), h = d;
        }
      } else
        for (let S in c)
          (y || c.hasOwnProperty(S)) && (d = h[S], d || (h[V] & 1048576 && (b = h[V] & 65535), d = h[S] = /* @__PURE__ */ Object.create(null), u++), h = d, x++);
      let A = h[V];
      if (A !== void 0)
        A &= 65535, a[i++] = 217, a[i++] = A >> 8 | 224, a[i++] = A & 255;
      else if (I || (I = h.__keys__ || (h.__keys__ = Object.keys(c))), b === void 0 ? (A = f.nextId++, A || (A = 0, f.nextId = 1), A >= Ht && (f.nextId = (A = C) + 1)) : A = b, f[A] = I, A < C) {
        a[i++] = 217, a[i++] = A >> 8 | 224, a[i++] = A & 255, h = f.transitions;
        for (let S = 0; S < x; S++)
          (h[V] === void 0 || h[V] & 1048576) && (h[V] = A), h = h[I[S]];
        h[V] = A | 1048576, o = true;
      } else {
        if (h[V] = A, M.setUint32(i, 3655335680), i += 3, u && (ve += De * u), J.length >= Ht - C && (J.shift()[V] = void 0), J.push(h), H(x + 2), U(57344 + A), U(I), y === null)
          return;
        for (let S in c)
          (y || c.hasOwnProperty(S)) && U(c[S]);
        return;
      }
      if (x < 24 ? a[i++] = 128 | x : H(x), y !== null)
        for (let S in c)
          (y || c.hasOwnProperty(S)) && U(c[S]);
    }, ue = (c) => {
      let y;
      if (c > 16777216) {
        if (c - n > qt)
          throw new Error("Encoded buffer would be larger than maximum buffer size");
        y = Math.min(qt, Math.round(Math.max((c - n) * (c > 67108864 ? 1.25 : 2), 4194304) / 4096) * 4096);
      } else
        y = (Math.max(c - n << 2, a.length - 1) >> 12) + 1 << 12;
      let d = new et(y);
      return M = new DataView(d.buffer, 0, y), a.copy ? a.copy(d, 0, n, c) : d.set(a.slice(n, c)), i -= n, n = 0, te = d.length - 10, a = d;
    }, se = 100, Ut = 1e3;
    this.encodeAsIterable = function(c, y) {
      return Ct(c, y, ae);
    }, this.encodeAsAsyncIterable = function(c, y) {
      return Ct(c, y, Mt);
    };
    function* ae(c, y, d) {
      let h = c.constructor;
      if (h === Object) {
        let u = w.useRecords !== false;
        u ? ze(c, null) : $t(Object.keys(c).length, 160);
        for (let x in c) {
          let b = c[x];
          u || U(x), b && typeof b == "object" ? y[x] ? yield* ae(b, y[x]) : yield* Ne(b, y, x) : U(b);
        }
      } else if (h === Array) {
        let u = c.length;
        H(u);
        for (let x = 0; x < u; x++) {
          let b = c[x];
          b && (typeof b == "object" || i - n > se) ? y.element ? yield* ae(b, y.element) : yield* Ne(b, y, "element") : U(b);
        }
      } else if (c[Symbol.iterator]) {
        a[i++] = 159;
        for (let u of c)
          u && (typeof u == "object" || i - n > se) ? y.element ? yield* ae(u, y.element) : yield* Ne(u, y, "element") : U(u);
        a[i++] = 255;
      } else
        nt(c) ? ($t(c.size, 64), yield a.subarray(n, i), yield c, he()) : c[Symbol.asyncIterator] ? (a[i++] = 159, yield a.subarray(n, i), yield c, he(), a[i++] = 255) : U(c);
      d && i > n ? yield a.subarray(n, i) : i - n > se && (yield a.subarray(n, i), he());
    }
    function* Ne(c, y, d) {
      let h = i - n;
      try {
        U(c), i - n > se && (yield a.subarray(n, i), he());
      } catch (u) {
        if (u.iteratorNotHandled)
          y[d] = {}, i = n + h, yield* ae.call(this, c, y[d]);
        else
          throw u;
      }
    }
    function he() {
      se = Ut, w.encode(null, at);
    }
    function Ct(c, y, d) {
      return y && y.chunkThreshold ? se = Ut = y.chunkThreshold : se = 100, c && typeof c == "object" ? (w.encode(null, at), d(c, w.iterateProperties || (w.iterateProperties = {}), true)) : [w.encode(c)];
    }
    async function* Mt(c, y) {
      for (let d of ae(c, y, true)) {
        let h = d.constructor;
        if (h === Kt || h === Uint8Array)
          yield d;
        else if (nt(d)) {
          let u = d.stream().getReader(), x;
          for (; !(x = await u.read()).done; )
            yield x.value;
        } else if (d[Symbol.asyncIterator])
          for await (let u of d)
            he(), u ? yield* Mt(u, y.async || (y.async = {})) : yield w.encode(u);
        else
          yield d;
      }
    }
  }
  useBuffer(e) {
    a = e, M = new DataView(a.buffer, a.byteOffset, a.byteLength), i = 0;
  }
  clearSharedData() {
    this.structures && (this.structures = []), this.sharedValues && (this.sharedValues = void 0);
  }
  updateSharedData() {
    let e = this.sharedVersion || 0;
    this.sharedVersion = e + 1;
    let r = this.structures.slice(0), n = new rt(r, this.sharedValues, this.sharedVersion), s = this.saveShared(n, (o) => (o && o.version || 0) == e);
    return s === false ? (n = this.getShared() || {}, this.structures = n.structures || [], this.sharedValues = n.packedValues, this.sharedVersion = n.version, this.structures.nextId = this.structures.length) : r.forEach((o, f) => this.structures[f] = o), s;
  }
};
function $t(t, e) {
  t < 24 ? a[i++] = e | t : t < 256 ? (a[i++] = e | 24, a[i++] = t) : t < 65536 ? (a[i++] = e | 25, a[i++] = t >> 8, a[i++] = t & 255) : (a[i++] = e | 26, M.setUint32(i, t), i += 4);
}
var rt = class {
  constructor(e, r, n) {
    this.structures = e, this.packedValues = r, this.version = n;
  }
};
function H(t) {
  t < 24 ? a[i++] = 128 | t : t < 256 ? (a[i++] = 152, a[i++] = t) : t < 65536 ? (a[i++] = 153, a[i++] = t >> 8, a[i++] = t & 255) : (a[i++] = 154, M.setUint32(i, t), i += 4);
}
var Mr = typeof Blob == "undefined" ? function() {
} : Blob;
function nt(t) {
  if (t instanceof Mr)
    return true;
  let e = t[Symbol.toStringTag];
  return e === "Blob" || e === "File";
}
function _e(t, e) {
  switch (typeof t) {
    case "string":
      if (t.length > 3) {
        if (e.objectMap[t] > -1 || e.values.length >= e.maxValues)
          return;
        let n = e.get(t);
        if (n)
          ++n.count == 2 && e.values.push(t);
        else if (e.set(t, { count: 1 }), e.samplingPackedValues) {
          let s = e.samplingPackedValues.get(t);
          s ? s.count++ : e.samplingPackedValues.set(t, { count: 1 });
        }
      }
      break;
    case "object":
      if (t)
        if (t instanceof Array)
          for (let n = 0, s = t.length; n < s; n++)
            _e(t[n], e);
        else {
          let n = !e.encoder.useRecords;
          for (var r in t)
            t.hasOwnProperty(r) && (n && _e(r, e), _e(t[r], e));
        }
      break;
    case "function":
      console.log(t);
  }
}
var Pr = new Uint8Array(new Uint16Array([1]).buffer)[0] == 1;
jt = [Date, Set, Error, RegExp, K, ArrayBuffer, Uint8Array, Uint8ClampedArray, Uint16Array, Uint32Array, typeof BigUint64Array == "undefined" ? function() {
} : BigUint64Array, Int8Array, Int16Array, Int32Array, typeof BigInt64Array == "undefined" ? function() {
} : BigInt64Array, Float32Array, Float64Array, rt];
Qe = [{ tag: 1, encode(t, e) {
  let r = t.getTime() / 1e3;
  (this.useTimestamp32 || t.getMilliseconds() === 0) && r >= 0 && r < 4294967296 ? (a[i++] = 26, M.setUint32(i, r), i += 4) : (a[i++] = 251, M.setFloat64(i, r), i += 8);
} }, { tag: 258, encode(t, e) {
  let r = Array.from(t);
  e(r);
} }, { tag: 27, encode(t, e) {
  e([t.name, t.message]);
} }, { tag: 27, encode(t, e) {
  e(["RegExp", t.source, t.flags]);
} }, { getTag(t) {
  return t.tag;
}, encode(t, e) {
  e(t.value);
} }, { encode(t, e, r) {
  Gt(t, r);
} }, { getTag(t) {
  if (t.constructor === Uint8Array && (this.tagUint8Array || me && this.tagUint8Array !== false))
    return 64;
}, encode(t, e, r) {
  Gt(t, r);
} }, q(68, 1), q(69, 2), q(70, 4), q(71, 8), q(72, 1), q(77, 2), q(78, 4), q(79, 8), q(85, 4), q(86, 8), { encode(t, e) {
  let r = t.packedValues || [], n = t.structures || [];
  if (r.values.length > 0) {
    a[i++] = 216, a[i++] = 51, H(4);
    let s = r.values;
    e(s), H(0), H(0), packedObjectMap = Object.create(sharedPackedObjectMap || null);
    for (let o = 0, f = s.length; o < f; o++)
      packedObjectMap[s[o]] = o;
  }
  if (n) {
    M.setUint32(i, 3655335424), i += 3;
    let s = n.slice(0);
    s.unshift(57344), s.push(new K(t.version, 1399353956)), e(s);
  } else
    e(new K(t.version, 1399353956));
} }];
function q(t, e) {
  return !Pr && e > 1 && (t -= 4), { tag: t, encode: function(n, s) {
    let o = n.byteLength, f = n.byteOffset || 0, m = n.buffer || n;
    s(me ? ke.from(m, f, o) : new Uint8Array(m, f, o));
  } };
}
function Gt(t, e) {
  let r = t.byteLength;
  r < 24 ? a[i++] = 64 + r : r < 256 ? (a[i++] = 88, a[i++] = r) : r < 65536 ? (a[i++] = 89, a[i++] = r >> 8, a[i++] = r & 255) : (a[i++] = 90, M.setUint32(i, r), i += 4), i + r >= a.length && e(i + r), a.set(t.buffer ? t : new Uint8Array(t), i), i += r;
}
function Or(t, e) {
  let r, n = e.length * 2, s = t.length - n;
  e.sort((o, f) => o.offset > f.offset ? 1 : -1);
  for (let o = 0; o < e.length; o++) {
    let f = e[o];
    f.id = o;
    for (let m of f.references)
      t[m++] = o >> 8, t[m] = o & 255;
  }
  for (; r = e.pop(); ) {
    let o = r.offset;
    t.copyWithin(o + n, o, s), n -= 2;
    let f = o + n;
    t[f++] = 216, t[f++] = 28, s = o;
  }
  return t;
}
function Jt(t, e) {
  M.setUint32(_.position + t, i - _.position - t + 1);
  let r = _;
  _ = null, e(r[0]), e(r[1]);
}
var st = new Re({ useRecords: false });
var it = st.encode;
var kr = st.encodeAsIterable;
var Rr = st.encodeAsAsyncIterable;
var ot = 512;
var _r = 1024;
var at = 2048;
var Br = class {
  constructor(e = false) {
    this.debug = e;
  }
  encoder(e) {
    return new Yt(e, this.debug);
  }
  decoder(e) {
    return new Xt(e, this.debug);
  }
};
var Yt = class {
  constructor(e, r = false) {
    this.w = e, this.debug = r;
  }
  async encode(e) {
    this.debug && console.log("<<", e);
    let r = it(e), n = 0;
    for (; n < r.length; )
      n += await this.w.write(r.subarray(n));
  }
};
var Xt = class {
  constructor(e, r = false) {
    this.r = e, this.debug = r;
  }
  async decode(e) {
    let r = new Uint8Array(e);
    if (await this.r.read(r) === null)
      return Promise.resolve(null);
    let s = Ze(r);
    return this.debug && console.log(">>", s), Promise.resolve(s);
  }
};
function Be(t, e, r = 0) {
  r = Math.max(0, Math.min(r, e.byteLength));
  let n = e.byteLength - r;
  return t.byteLength > n && (t = t.subarray(0, n)), e.set(t, r), t.byteLength;
}
var Fe = 32 * 1024;
var ct = 2 ** 32 - 2;
var lt = class {
  constructor(e) {
    this._buf = e === void 0 ? new Uint8Array(0) : new Uint8Array(e), this._off = 0;
  }
  bytes(e = { copy: true }) {
    return e.copy === false ? this._buf.subarray(this._off) : this._buf.slice(this._off);
  }
  empty() {
    return this._buf.byteLength <= this._off;
  }
  get length() {
    return this._buf.byteLength - this._off;
  }
  get capacity() {
    return this._buf.buffer.byteLength;
  }
  truncate(e) {
    if (e === 0) {
      this.reset();
      return;
    }
    if (e < 0 || e > this.length)
      throw Error("bytes.Buffer: truncation out of range");
    this._reslice(this._off + e);
  }
  reset() {
    this._reslice(0), this._off = 0;
  }
  _tryGrowByReslice(e) {
    let r = this._buf.byteLength;
    return e <= this.capacity - r ? (this._reslice(r + e), r) : -1;
  }
  _reslice(e) {
    this._buf = new Uint8Array(this._buf.buffer, 0, e);
  }
  readSync(e) {
    if (this.empty())
      return this.reset(), e.byteLength === 0 ? 0 : null;
    let r = Be(this._buf.subarray(this._off), e);
    return this._off += r, r;
  }
  read(e) {
    let r = this.readSync(e);
    return Promise.resolve(r);
  }
  writeSync(e) {
    let r = this._grow(e.byteLength);
    return Be(e, this._buf, r);
  }
  write(e) {
    let r = this.writeSync(e);
    return Promise.resolve(r);
  }
  _grow(e) {
    let r = this.length;
    r === 0 && this._off !== 0 && this.reset();
    let n = this._tryGrowByReslice(e);
    if (n >= 0)
      return n;
    let s = this.capacity;
    if (e <= Math.floor(s / 2) - r)
      Be(this._buf.subarray(this._off), this._buf);
    else {
      if (s + e > ct)
        throw new Error("The buffer cannot be grown beyond the maximum size.");
      {
        let o = new Uint8Array(Math.min(2 * s + e, ct));
        Be(this._buf.subarray(this._off), o), this._buf = o;
      }
    }
    return this._off = 0, this._reslice(Math.min(r + e, ct)), r;
  }
  grow(e) {
    if (e < 0)
      throw Error("Buffer.grow: negative count");
    let r = this._grow(e);
    this._reslice(r);
  }
  async readFrom(e) {
    let r = 0, n = new Uint8Array(Fe);
    for (; ; ) {
      let s = this.capacity - this.length < Fe, o = s ? n : new Uint8Array(this._buf.buffer, this.length), f = await e.read(o);
      if (f === null)
        return r;
      s ? this.writeSync(o.subarray(0, f)) : this._reslice(this.length + f), r += f;
    }
  }
  readFromSync(e) {
    let r = 0, n = new Uint8Array(Fe);
    for (; ; ) {
      let s = this.capacity - this.length < Fe, o = s ? n : new Uint8Array(this._buf.buffer, this.length), f = e.readSync(o);
      if (f === null)
        return r;
      s ? this.writeSync(o.subarray(0, f)) : this._reslice(this.length + f), r += f;
    }
  }
};
var xe = class {
  constructor(e) {
    this.codec = e;
  }
  encoder(e) {
    return new Zt(e, this.codec);
  }
  decoder(e) {
    return new Qt(e, this.codec.decoder(e));
  }
};
var Zt = class {
  constructor(e, r) {
    this.w = e, this.codec = r;
  }
  async encode(e) {
    let r = new lt();
    await this.codec.encoder(r).encode(e);
    let s = new DataView(new ArrayBuffer(4));
    s.setUint32(0, r.length);
    let o = new Uint8Array(r.length + 4);
    o.set(new Uint8Array(s.buffer), 0), o.set(r.bytes(), 4);
    let f = 0;
    for (; f < o.length; )
      f += await this.w.write(o.subarray(f));
  }
};
var Qt = class {
  constructor(e, r) {
    this.r = e, this.dec = r;
  }
  async decode(e) {
    let r = new Uint8Array(4);
    if (await this.r.read(r) === null)
      return null;
    let o = new DataView(r.buffer).getUint32(0);
    return await this.dec.decode(o);
  }
};
var ge = class {
  constructor(e, r) {
    this.session = e, this.codec = r;
  }
  async call(e, r) {
    let n = await this.session.open();
    try {
      let s = new xe(this.codec), o = s.encoder(n), f = s.decoder(n);
      await o.encode({ Selector: e }), await o.encode(r);
      let m = await f.decode(), g = new ft(n, s);
      if (g.error = m.Error, g.error !== void 0 && g.error !== null)
        throw g.error;
      return g.reply = await f.decode(), g.continue = m.Continue, g.continue || await n.close(), g;
    } catch (s) {
      return await n.close(), console.error(s, e, r), Promise.reject(s);
    }
  }
};
function tr(t) {
  function e(r, n) {
    return new Proxy(Object.assign(() => {
    }, { path: r, callable: n }), { get(s, o, f) {
      return o.startsWith("__") ? Reflect.get(s, o, f) : e(s.path ? `${s.path}.${o}` : o, s.callable);
    }, apply(s, o, f = []) {
      return s.callable(s.path, f);
    } });
  }
  return e("", t.call.bind(t));
}
function dt(t) {
  return { respondRPC: t };
}
function Fr() {
  return dt((t, e) => {
    t.return(new Error(`not found: ${e.selector}`));
  });
}
function ut(t) {
  return t === "" ? "/" : (t[0] != "/" && (t = "/" + t), t = t.replace(".", "/"), t);
}
var be = class {
  constructor() {
    this.handlers = {};
  }
  async respondRPC(e, r) {
    await this.handler(r).respondRPC(e, r);
  }
  handler(e) {
    let r = this.match(e.selector);
    return r || Fr();
  }
  remove(e) {
    e = ut(e);
    let r = this.match(e);
    return delete this.handlers[e], r || null;
  }
  match(e) {
    return e = ut(e), this.handlers.hasOwnProperty(e) ? this.handlers[e] : null;
  }
  handle(e, r) {
    if (e === "")
      throw "invalid selector";
    if (e = ut(e), !r)
      throw "invalid handler";
    if (this.match(e))
      throw "selector already registered";
    this.handlers[e] = r;
  }
};
async function rr(t, e, r) {
  let n = new xe(e), s = n.decoder(t), o = await s.decode(), f = new ht(o.Selector, s);
  f.caller = new ge(t.session, e);
  let m = new yt(), g = new nr(t, n, m);
  return r || (r = new be()), await r.respondRPC(g, f), g.responded || await g.return(null), Promise.resolve();
}
var nr = class {
  constructor(e, r, n) {
    this.ch = e, this.codec = r, this.header = n, this.responded = false;
  }
  send(e) {
    return this.codec.encoder(this.ch).encode(e);
  }
  return(e) {
    return this.respond(e, false);
  }
  async continue(e) {
    return await this.respond(e, true), this.ch;
  }
  async respond(e, r) {
    return this.responded = true, this.header.Continue = r, e instanceof Error && (this.header.Error = e.message, e = null), await this.send(this.header), await this.send(e), r || await this.ch.close(), Promise.resolve();
  }
};
var ht = class {
  constructor(e, r) {
    this.selector = e, this.decoder = r;
  }
  receive() {
    return this.decoder.decode();
  }
};
var yt = class {
  constructor() {
    this.Error = void 0, this.Continue = false;
  }
};
var ft = class {
  constructor(e, r) {
    this.channel = e, this.codec = r, this.error = void 0, this.continue = false;
  }
  send(e) {
    this.codec.encoder(this.channel).encode(e);
  }
  receive() {
    return this.codec.decoder(this.channel).decode();
  }
};
var pt = class {
  constructor(e, r) {
    this.session = e, this.codec = r, this.caller = new ge(e, r), this.responder = new be();
  }
  async respond() {
    for (; ; ) {
      let e = await this.session.accept();
      if (e === null)
        break;
      rr(e, this.codec, this.responder);
    }
  }
  async call(e, r) {
    return this.caller.call(e, r);
  }
  handle(e, r) {
    this.responder.handle(e, r);
  }
  respondRPC(e, r) {
    this.responder.respondRPC(e, r);
  }
  virtualize() {
    return tr(this.caller);
  }
};
var X = 100;
var Z = 101;
var Q = 102;
var re = 103;
var ee = 104;
var ne = 105;
var $ = 106;
var sr = /* @__PURE__ */ new Map([[X, 12], [Z, 16], [Q, 4], [re, 8], [ee, 8], [ne, 4], [$, 4]]);
var wt = class {
  constructor(e) {
    this.w = e;
  }
  async encode(e) {
    de.messages && console.log("<<ENC", e);
    let r = Tr(e);
    de.bytes && console.log("<<ENC", r);
    let n = 0;
    for (; n < r.length; )
      n += await this.w.write(r.subarray(n));
    return n;
  }
};
function Tr(t) {
  if (t.ID === $) {
    let e = t, r = new DataView(new ArrayBuffer(5));
    return r.setUint8(0, e.ID), r.setUint32(1, e.channelID), new Uint8Array(r.buffer);
  }
  if (t.ID === ee) {
    let e = t, r = new DataView(new ArrayBuffer(9));
    r.setUint8(0, e.ID), r.setUint32(1, e.channelID), r.setUint32(5, e.length);
    let n = new Uint8Array(9 + e.length);
    return n.set(new Uint8Array(r.buffer), 0), n.set(e.data, 9), n;
  }
  if (t.ID === ne) {
    let e = t, r = new DataView(new ArrayBuffer(5));
    return r.setUint8(0, e.ID), r.setUint32(1, e.channelID), new Uint8Array(r.buffer);
  }
  if (t.ID === X) {
    let e = t, r = new DataView(new ArrayBuffer(13));
    return r.setUint8(0, e.ID), r.setUint32(1, e.senderID), r.setUint32(5, e.windowSize), r.setUint32(9, e.maxPacketSize), new Uint8Array(r.buffer);
  }
  if (t.ID === Z) {
    let e = t, r = new DataView(new ArrayBuffer(17));
    return r.setUint8(0, e.ID), r.setUint32(1, e.channelID), r.setUint32(5, e.senderID), r.setUint32(9, e.windowSize), r.setUint32(13, e.maxPacketSize), new Uint8Array(r.buffer);
  }
  if (t.ID === Q) {
    let e = t, r = new DataView(new ArrayBuffer(5));
    return r.setUint8(0, e.ID), r.setUint32(1, e.channelID), new Uint8Array(r.buffer);
  }
  if (t.ID === re) {
    let e = t, r = new DataView(new ArrayBuffer(9));
    return r.setUint8(0, e.ID), r.setUint32(1, e.channelID), r.setUint32(5, e.additionalBytes), new Uint8Array(r.buffer);
  }
  throw `marshal of unknown type: ${t}`;
}
function Ve(t, e) {
  let r = new Uint8Array(e), n = 0;
  return t.forEach((s) => {
    r.set(s, n), n += s.length;
  }), r;
}
var Ie = class {
  constructor() {
    this.q = [], this.waiters = [], this.closed = false;
  }
  push(e) {
    if (this.closed)
      throw "closed queue";
    if (this.waiters.length > 0) {
      let r = this.waiters.shift();
      r && r(e);
      return;
    }
    this.q.push(e);
  }
  shift() {
    return this.closed ? Promise.resolve(null) : new Promise((e) => {
      if (this.q.length > 0) {
        e(this.q.shift() || null);
        return;
      }
      this.waiters.push(e);
    });
  }
  close() {
    this.closed || (this.closed = true, this.waiters.forEach((e) => {
      e(null);
    }));
  }
};
var mt = class {
  constructor() {
    this.readBuf = new Uint8Array(0), this.gotEOF = false, this.readers = [];
  }
  read(e) {
    return new Promise((r) => {
      let n = () => {
        if (this.readBuf === void 0) {
          r(null);
          return;
        }
        if (this.readBuf.length == 0) {
          if (this.gotEOF) {
            this.readBuf = void 0, r(null);
            return;
          }
          this.readers.push(n);
          return;
        }
        let s = this.readBuf.slice(0, e.length);
        this.readBuf = this.readBuf.slice(s.length), this.readBuf.length == 0 && this.gotEOF && (this.readBuf = void 0), e.set(s), r(s.length);
      };
      n();
    });
  }
  write(e) {
    for (this.readBuf && (this.readBuf = Ve([this.readBuf, e], this.readBuf.length + e.length)); !this.readBuf || this.readBuf.length > 0; ) {
      let r = this.readers.shift();
      if (!r)
        break;
      r();
    }
    return Promise.resolve(e.length);
  }
  eof() {
    this.gotEOF = true, this.flushReaders();
  }
  close() {
    this.readBuf = void 0, this.flushReaders();
  }
  flushReaders() {
    for (; ; ) {
      let e = this.readers.shift();
      if (!e)
        return;
      e();
    }
  }
};
var gt = class {
  constructor(e) {
    this.r = e;
  }
  async decode() {
    let e = await Vr(this.r);
    if (e === null)
      return Promise.resolve(null);
    de.bytes && console.log(">>DEC", e);
    let r = Lr(e);
    return de.messages && console.log(">>DEC", r), r;
  }
};
async function Vr(t) {
  let e = new Uint8Array(1);
  if (await t.read(e) === null)
    return Promise.resolve(null);
  let n = e[0], s = sr.get(n);
  if (s === void 0 || n < X || n > $)
    return Promise.reject(`bad packet: ${n}`);
  let o = new Uint8Array(s);
  if (await t.read(o) === null)
    return Promise.reject("unexpected EOF");
  if (n === ee) {
    let g = new DataView(o.buffer).getUint32(4), w = new Uint8Array(g);
    return await t.read(w) === null ? Promise.reject("unexpected EOF") : Ve([e, o, w], g + o.length + 1);
  }
  return Ve([e, o], o.length + 1);
}
function Lr(t) {
  let e = new DataView(t.buffer);
  switch (t[0]) {
    case $:
      return { ID: t[0], channelID: e.getUint32(1) };
    case ee:
      let r = e.getUint32(5), n = new Uint8Array(t.buffer.slice(9));
      return { ID: t[0], channelID: e.getUint32(1), length: r, data: n };
    case ne:
      return { ID: t[0], channelID: e.getUint32(1) };
    case X:
      return { ID: t[0], senderID: e.getUint32(1), windowSize: e.getUint32(5), maxPacketSize: e.getUint32(9) };
    case Z:
      return { ID: t[0], channelID: e.getUint32(1), senderID: e.getUint32(5), windowSize: e.getUint32(9), maxPacketSize: e.getUint32(13) };
    case Q:
      return { ID: t[0], channelID: e.getUint32(1) };
    case re:
      return { ID: t[0], channelID: e.getUint32(1), additionalBytes: e.getUint32(5) };
    default:
      throw `unmarshal of unknown type: ${t[0]}`;
  }
}
var de = { messages: false, bytes: false };
var bt = 9;
var It = Number.MAX_VALUE;
var At = class {
  constructor(e) {
    this.conn = e, this.enc = new wt(e), this.dec = new gt(e), this.channels = [], this.incoming = new Ie(), this.done = this.loop();
  }
  async open() {
    let e = this.newChannel();
    if (e.maxIncomingPayload = Le, await this.enc.encode({ ID: X, windowSize: e.myWindow, maxPacketSize: e.maxIncomingPayload, senderID: e.localId }), await e.ready.shift())
      return e;
    throw "failed to open";
  }
  accept() {
    return this.incoming.shift();
  }
  async close() {
    for (let e of Object.keys(this.channels)) {
      let r = parseInt(e);
      this.channels[r] !== void 0 && this.channels[r].shutdown();
    }
    this.conn.close(), await this.done;
  }
  async loop() {
    try {
      for (; ; ) {
        let e = await this.dec.decode();
        if (e === null) {
          this.close();
          return;
        }
        if (e.ID === X) {
          await this.handleOpen(e);
          continue;
        }
        let r = e, n = this.getCh(r.channelID);
        if (n === void 0)
          throw `invalid channel (${r.channelID}) on op ${r.ID}`;
        await n.handle(r);
      }
    } catch (e) {
      throw new Error(`session loop: ${e}`);
    }
  }
  async handleOpen(e) {
    if (e.maxPacketSize < bt || e.maxPacketSize > It) {
      await this.enc.encode({ ID: Q, channelID: e.senderID });
      return;
    }
    let r = this.newChannel();
    r.remoteId = e.senderID, r.maxRemotePayload = e.maxPacketSize, r.remoteWin = e.windowSize, r.maxIncomingPayload = Le, this.incoming.push(r), await this.enc.encode({ ID: Z, channelID: r.remoteId, senderID: r.localId, windowSize: r.myWindow, maxPacketSize: r.maxIncomingPayload });
  }
  newChannel() {
    let e = new Dt(this);
    return e.remoteWin = 0, e.myWindow = ir, e.localId = this.addCh(e), e;
  }
  getCh(e) {
    let r = this.channels[e];
    return r && r.localId !== e && console.log("bad ids:", e, r.localId, r.remoteId), r;
  }
  addCh(e) {
    return this.channels.forEach((r, n) => {
      if (r === void 0)
        return this.channels[n] = e, n;
    }), this.channels.push(e), this.channels.length - 1;
  }
  rmCh(e) {
    delete this.channels[e];
  }
};
var Le = 1 << 24;
var ir = 64 * Le;
var Dt = class {
  constructor(e) {
    this.localId = 0, this.remoteId = 0, this.maxIncomingPayload = 0, this.maxRemotePayload = 0, this.sentEOF = false, this.sentClose = false, this.remoteWin = 0, this.myWindow = 0, this.ready = new Ie(), this.session = e, this.writers = [], this.readBuf = new mt();
  }
  ident() {
    return this.localId;
  }
  async read(e) {
    let r = await this.readBuf.read(e);
    if (r !== null)
      try {
        await this.adjustWindow(r);
      } catch (n) {
        if (n !== "EOF")
          throw n;
      }
    return r;
  }
  write(e) {
    return this.sentEOF ? Promise.reject("EOF") : new Promise((r, n) => {
      let s = 0, o = () => {
        if (this.sentEOF || this.sentClose) {
          n("EOF");
          return;
        }
        if (e.byteLength == 0) {
          r(s);
          return;
        }
        let f = Math.min(this.maxRemotePayload, e.length), m = this.reserveWindow(f);
        if (m == 0) {
          this.writers.push(o);
          return;
        }
        let g = e.slice(0, m);
        this.send({ ID: ee, channelID: this.remoteId, length: g.length, data: g }).then(() => {
          if (s += g.length, e = e.slice(g.length), e.length == 0) {
            r(s);
            return;
          }
          this.writers.push(o);
        });
      };
      o();
    });
  }
  reserveWindow(e) {
    return this.remoteWin < e && (e = this.remoteWin), this.remoteWin -= e, e;
  }
  addWindow(e) {
    for (this.remoteWin += e; this.remoteWin > 0; ) {
      let r = this.writers.shift();
      if (!r)
        break;
      r();
    }
  }
  async closeWrite() {
    this.sentEOF = true, await this.send({ ID: ne, channelID: this.remoteId }), this.writers.forEach((e) => e()), this.writers = [];
  }
  async close() {
    if (!this.sentClose) {
      for (await this.send({ ID: $, channelID: this.remoteId }), this.sentClose = true; await this.ready.shift() !== null; )
        ;
      return;
    }
    this.shutdown();
  }
  shutdown() {
    this.readBuf.close(), this.writers.forEach((e) => e()), this.ready.close(), this.session.rmCh(this.localId);
  }
  async adjustWindow(e) {
    this.myWindow += e, await this.send({ ID: re, channelID: this.remoteId, additionalBytes: e });
  }
  send(e) {
    if (this.sentClose)
      throw "EOF";
    return this.sentClose = e.ID === $, this.session.enc.encode(e);
  }
  handle(e) {
    if (e.ID === ee) {
      this.handleData(e);
      return;
    }
    if (e.ID === $) {
      this.close();
      return;
    }
    if (e.ID === ne && this.readBuf.eof(), e.ID === Q) {
      this.session.rmCh(e.channelID), this.ready.push(false);
      return;
    }
    if (e.ID === Z) {
      if (e.maxPacketSize < bt || e.maxPacketSize > It)
        throw "invalid max packet size";
      this.remoteId = e.senderID, this.maxRemotePayload = e.maxPacketSize, this.addWindow(e.windowSize), this.ready.push(true);
      return;
    }
    e.ID === re && this.addWindow(e.additionalBytes);
  }
  handleData(e) {
    if (e.length > this.maxIncomingPayload)
      throw "incoming packet exceeds maximum payload size";
    if (this.myWindow < e.length)
      throw "remote side wrote too much";
    this.myWindow -= e.length, this.readBuf.write(e.data);
  }
};
var St = {};
fr(St, { Conn: () => Et, connect: () => zr });
function zr(t, e) {
  return new Promise((r) => {
    let n = new WebSocket(t);
    n.onopen = () => r(new Et(n)), e && (n.onclose = e);
  });
}
var Et = class {
  constructor(e) {
    this.isClosed = false, this.waiters = [], this.chunks = [], this.ws = e, this.ws.binaryType = "arraybuffer", this.ws.onmessage = (n) => {
      let s = new Uint8Array(n.data);
      if (this.chunks.push(s), this.waiters.length > 0) {
        let o = this.waiters.shift();
        o && o();
      }
    };
    let r = this.ws.onclose;
    this.ws.onclose = (n) => {
      r && r.bind(this.ws)(n), this.close();
    };
  }
  read(e) {
    return new Promise((r) => {
      var n = () => {
        if (this.isClosed) {
          r(null);
          return;
        }
        if (this.chunks.length === 0) {
          this.waiters.push(n);
          return;
        }
        let s = 0;
        for (; s < e.length; ) {
          let o = this.chunks.shift();
          if (o == null) {
            r(null);
            return;
          }
          let f = o.slice(0, e.length - s);
          if (e.set(f, s), s += f.length, o.length > f.length) {
            let m = o.slice(f.length);
            this.chunks.unshift(m);
          }
        }
        r(s);
      };
      n();
    });
  }
  write(e) {
    return this.ws.send(e), Promise.resolve(e.byteLength);
  }
  close() {
    this.isClosed || (this.isClosed = true, this.waiters.forEach((e) => e()), this.ws.close());
  }
};
var Nr = { transport: St };
async function Nn(t, e) {
  let r = await Nr.transport.connect(t);
  return jr(r, e);
}
function jr(t, e, r) {
  let n = new At(t), s = new pt(n, e);
  if (r) {
    for (let o in r)
      s.handle(o, dt(r[o]));
    s.respond();
  }
  return s;
}
var Ae = {};
function ar(t) {
  return t.frameElement ? t.frameElement.id : "";
}
window.addEventListener("message", (t) => {
  if (!t.source)
    return;
  let e = ar(t.source);
  if (!Ae[e]) {
    let s = new CustomEvent("connection", { detail: e });
    if (!window.dispatchEvent(s))
      return;
    if (!Ae[e]) {
      console.warn("incoming message with no connection for frame ID in window:", e, window.location);
      return;
    }
  }
  let r = Ae[e], n = new Uint8Array(t.data);
  if (r.chunks.push(n), r.waiters.length > 0) {
    let s = r.waiters.shift();
    s && s();
  }
});

// src/client.ts
async function connect(url) {
  return new Client(await Nn(url, new Br()));
}
var Client = class {
  constructor(peer) {
    this.data = {};
    this.peer = peer;
    this.rpc = peer.virtualize();
    this.app = new AppModule(this);
    this.menu = new MenuModule(this);
    this.system = new SystemModule(this);
    this.shell = new ShellModule(this);
    this.window = new WindowModule(this);
    this.handleEvents(peer);
    this.peer.respond();
  }
  async handleEvents(peer) {
    const resp = await peer.call("Listen");
    while (true) {
      const obj = await resp.receive();
      if (obj === null) {
        break;
      }
      const event = obj;
      if (this.onevent) {
        this.onevent(event);
      }
      switch (event.Name) {
        case "menu":
          if (this.menu.onclick)
            this.menu.onclick(event);
          break;
        case "shortcut":
          if (this.shell.onshortcut)
            this.shell.onshortcut(event);
          break;
        default:
          const w = this.window.windows[event.WindowID];
          if (w) {
            switch (event.Name) {
              case "close":
                if (w.onclose)
                  w.onclose(event);
                break;
              case "destroy":
                if (w.ondestroyed)
                  w.ondestroyed(event);
                delete this.window.windows[event.WindowID];
                break;
              case "focus":
                if (w.onfocused)
                  w.onfocused(event);
                break;
              case "blur":
                if (w.onblurred)
                  w.onblurred(event);
                break;
              case "resize":
                if (w.onresized)
                  w.onresized(event);
                break;
              case "move":
                if (w.onmoved)
                  w.onmoved(event);
                break;
            }
          }
      }
    }
  }
};
var AppModule = class {
  constructor(client) {
    this.rpc = client.rpc;
    this.client = client;
  }
  Run(options) {
    this.rpc.app.Run(options);
  }
  Menu() {
    return this.rpc.app.Menu();
  }
  SetMenu(m) {
    this.rpc.app.SetMenu(m);
  }
  async NewIndicator(icon, items) {
    this.rpc.app.NewIndicator(icon, items);
  }
};
var MenuModule = class {
  constructor(client) {
    this.rpc = client.rpc;
  }
  New(items) {
    return this.rpc.menu.New(items);
  }
  Popup(items) {
    return this.rpc.menu.Popup(items);
  }
};
var SystemModule = class {
  constructor(client) {
    this.rpc = client.rpc;
  }
  Displays() {
    return this.rpc.system.Displays();
  }
};
var ShellModule = class {
  constructor(client) {
    this.rpc = client.rpc;
  }
  ShowNotification(n) {
    this.rpc.shell.ShowNotification(n);
  }
  ShowMessage(msg) {
    this.rpc.shell.ShowMessage(msg);
  }
  ShowFilePicker(fd) {
    return this.rpc.shell.ShowFilePicker(fd);
  }
  ReadClipboard() {
    return this.rpc.shell.ReadClipboard();
  }
  WriteClipboard(text) {
    return this.rpc.shell.WriteClipboard(text);
  }
  RegisterShortcut(accelerator) {
    return this.rpc.shell.RegisterShortcut(accelerator);
  }
  IsShortcutRegistered(accelerator) {
    return this.rpc.shell.IsShortcutRegistered(accelerator);
  }
  UnregisterShortcut(accelerator) {
    return this.rpc.shell.UnregisterShortcut(accelerator);
  }
  UnregisterAllShortcuts() {
    return this.rpc.shell.UnregisterAllShortcuts();
  }
};
var WindowModule = class {
  constructor(client) {
    this.rpc = client.rpc;
    this.main = new Window(this.rpc, 0);
    this.windows = { 0: this.main };
  }
  async New(options) {
    const w = await this.rpc.window.New(options);
    this.windows[w.ID] = new Window(this.rpc, w.ID);
    return this.windows[w.ID];
  }
};
var Menu = class {
  constructor(rpc, id) {
    this.rpc = rpc;
    this.ID = id;
  }
};
var Window = class {
  constructor(rpc, id) {
    this.rpc = rpc;
    this.ID = id;
  }
  destroy() {
    this.rpc.window.Destroy(this.ID);
  }
  focus() {
    this.rpc.window.Focus(this.ID);
  }
  getOuterPosition() {
    return this.rpc.window.GetOuterPosition(this.ID);
  }
  getOuterSize() {
    return this.rpc.window.GetOuterSize(this.ID);
  }
  isDestroyed() {
    return this.rpc.window.IsDestroyed(this.ID);
  }
  isVisible() {
    return this.rpc.window.IsVisible(this.ID);
  }
  setVisible(visible) {
    return this.rpc.window.SetVisible(this.ID, visible);
  }
  setMaximized(maximized) {
    return this.rpc.window.SetMaximized(this.ID, maximized);
  }
  setMinimized(minimized) {
    return this.rpc.window.SetMinimized(this.ID, minimized);
  }
  setFullscreen(fullscreen) {
    return this.rpc.window.SetFullscreen(this.ID, fullscreen);
  }
  setMinSize(size) {
    return this.rpc.window.SetMinSize(this.ID, size);
  }
  setMaxSize(size) {
    return this.rpc.window.SetMaxSize(this.ID, size);
  }
  setResizable(resizable) {
    return this.rpc.window.SetResizable(this.ID, resizable);
  }
  setAlwaysOnTop(always) {
    return this.rpc.window.SetAlwaysOnTop(this.ID, always);
  }
  setSize(size) {
    return this.rpc.window.SetSize(this.ID, size);
  }
  setPosition(position) {
    return this.rpc.window.SetPosition(this.ID, position);
  }
  setTitle(title) {
    return this.rpc.window.SetTitle(this.ID, title);
  }
};
export {
  Client,
  Menu,
  Window,
  connect
};

package com.oceancatch;

import android.content.Context;
import android.graphics.Canvas;
import android.graphics.Color;
import android.graphics.Paint;
import android.graphics.RectF;
import android.os.SystemClock;
import android.view.MotionEvent;
import android.view.View;

import java.util.ArrayList;
import java.util.List;
import java.util.Random;

public class OceanView extends View {
    private final Paint p = new Paint(Paint.ANTI_ALIAS_FLAG);
    private final Random random = new Random();
    private final List<Fish> fish = new ArrayList<>();
    private long lastFrame = 0L;
    private int screen = 0;
    private int spot = 0;
    private int coins = 100;
    private int caught = 0;
    private int selectedFish = -1;
    private final boolean[] codex = new boolean[12];
    private float density;

    private static final String[] NAMES = {
            "고등어", "참돔", "광어", "오징어", "방어", "연어",
            "참치", "청새치", "문어", "복어", "상어", "만타"
    };
    private static final int[] VALUES = {8, 12, 15, 18, 25, 30, 45, 70, 40, 22, 95, 120};

    public OceanView(Context context) {
        super(context);
        density = getResources().getDisplayMetrics().density;
        setFocusable(true);
        for (int i = 0; i < 22; i++) fish.add(new Fish(i % 12, random.nextFloat(), 0.28f + random.nextFloat() * 0.55f));
    }

    private float s(float v) { return v * density; }

    @Override protected void onDraw(Canvas c) {
        super.onDraw(c);
        float w = getWidth(), h = getHeight();
        drawBackground(c, w, h);
        if (screen == 0) drawFishing(c, w, h);
        else if (screen == 1) drawMarket(c, w, h);
        else drawCodex(c, w, h);
        invalidate();
    }

    private void drawBackground(Canvas c, float w, float h) {
        p.setStyle(Paint.Style.FILL);
        p.setColor(Color.rgb(5, 25, 45)); c.drawRect(0, 0, w, h, p);
        p.setColor(Color.rgb(8, 62, 98)); c.drawRect(0, h * .30f, w, h, p);
        p.setColor(Color.rgb(9, 83, 120)); c.drawRect(0, h * .48f, w, h, p);
        p.setColor(Color.rgb(30, 150, 190));
        for (int i = 0; i < 7; i++) {
            float y = h * (.33f + i * .09f);
            c.drawRect(0, y, w, y + s(2), p);
        }
        p.setColor(Color.rgb(245, 192, 84)); c.drawCircle(w - s(55), s(55), s(23), p);
        text(c, "OCEAN CATCH", s(20), s(34), s(24), Color.WHITE, true);
        text(c, "💰 " + coins + "G   🐟 " + caught + "마리", s(20), s(66), s(16), Color.rgb(230, 245, 250), false);
    }

    private void drawFishing(Canvas c, float w, float h) {
        text(c, "낚시터", s(20), h * .28f, s(23), Color.WHITE, true);
        String[] spots = {"연안", "깊은 바다", "푸른 해구"};
        for (int i = 0; i < 3; i++) {
            float x = s(18 + i * 150);
            button(c, x, h * .30f, s(136), s(48), spots[i], i == spot);
        }

        long now = SystemClock.uptimeMillis();
        if (lastFrame == 0) lastFrame = now;
        float dt = Math.min(0.05f, (now - lastFrame) / 1000f); lastFrame = now;
        for (int i = 0; i < fish.size(); i++) {
            Fish f = fish.get(i);
            f.x += f.speed * dt * (0.55f + spot * .12f);
            if (f.x > 1.12f) f.x = -0.12f;
            float yy = h * (.58f + (i % 5) * .075f);
            drawFish(c, f.x * w, yy, s(13 + (f.kind % 3) * 3), f.kind);
        }

        if (selectedFish >= 0) {
            Fish f = fish.get(selectedFish);
            p.setColor(Color.argb(235, 7, 20, 33));
            c.drawRoundRect(new RectF(s(28), h * .58f, w - s(28), h - s(55)), s(18), s(18), p);
            text(c, "입질 성공!", s(52), h * .67f, s(25), Color.rgb(125, 220, 255), true);
            text(c, NAMES[f.kind] + "  +" + VALUES[f.kind] + "G", s(52), h * .73f, s(21), Color.WHITE, true);
            text(c, "수조에 보관하거나 시장에 판매할 수 있어요.", s(52), h * .78f, s(14), Color.LTGRAY, false);
            button(c, s(52), h * .83f, s(120), s(46), "수조", false);
            button(c, s(188), h * .83f, s(120), s(46), "즉시 판매", false);
        } else {
            button(c, w * .5f - s(90), h - s(70), s(180), s(50), "🎣 낚시하기", false);
        }
        nav(c, w, h);
    }

    private void drawMarket(Canvas c, float w, float h) {
        text(c, "수산시장", s(20), s(105), s(24), Color.WHITE, true);
        text(c, "보유한 물고기를 판매하면 골드를 얻습니다.", s(20), s(132), s(14), Color.LTGRAY, false);
        for (int i = 0; i < 6; i++) {
            float x = s(18 + (i % 2) * 225), y = s(160 + (i / 2) * 92);
            p.setColor(Color.argb(170, 255, 255, 255));
            c.drawRoundRect(new RectF(x, y, x + s(205), y + s(76)), s(15), s(15), p);
            int k = (i + spot) % NAMES.length;
            drawFish(c, x + s(30), y + s(38), s(16), k);
            text(c, NAMES[k], x + s(56), y + s(29), s(17), Color.WHITE, true);
            text(c, VALUES[k] + "G", x + s(56), y + s(53), s(14), Color.rgb(255, 214, 95), true);
            button(c, x + s(128), y + s(16), s(65), s(44), "판매", false);
        }
        nav(c, w, h);
    }

    private void drawCodex(Canvas c, float w, float h) {
        text(c, "도감", s(20), s(105), s(24), Color.WHITE, true);
        text(c, "발견한 어종은 자동으로 기록됩니다.", s(20), s(132), s(14), Color.LTGRAY, false);
        for (int i = 0; i < 12; i++) {
            float x = s(18 + (i % 3) * 150), y = s(158 + (i / 3) * 92);
            p.setColor(Color.argb(180, 255, 255, 255));
            c.drawRoundRect(new RectF(x, y, x + s(134), y + s(78)), s(15), s(15), p);
            if (codex[i]) {
                drawFish(c, x + s(26), y + s(38), s(14), i);
                text(c, NAMES[i], x + s(47), y + s(35), s(15), Color.WHITE, true);
                text(c, "판매가 " + VALUES[i] + "G", x + s(47), y + s(58), s(12), Color.LTGRAY, false);
            } else {
                text(c, "???", x + s(47), y + s(39), s(19), Color.LTGRAY, true);
            }
        }
        nav(c, w, h);
    }

    private void nav(Canvas c, float w, float h) {
        float y = h - s(62);
        String[] labels = {"🎣 낚시", "🧺 시장", "📖 도감"};
        for (int i = 0; i < 3; i++) button(c, s(10 + i * 112), y, s(102), s(48), labels[i], screen == i);
    }

    private void button(Canvas c, float x, float y, float ww, float hh, String label, boolean selected) {
        p.setColor(selected ? Color.rgb(30, 151, 205) : Color.rgb(19, 48, 70));
        c.drawRoundRect(new RectF(x, y, x + ww, y + hh), s(13), s(13), p);
        p.setStyle(Paint.Style.STROKE); p.setStrokeWidth(s(1)); p.setColor(Color.argb(110, 180, 235, 255));
        c.drawRoundRect(new RectF(x, y, x + ww, y + hh), s(13), s(13), p); p.setStyle(Paint.Style.FILL);
        text(c, label, x + ww / 2, y + hh * .63f, s(15), Color.WHITE, true, true);
    }

    private void drawFish(Canvas c, float x, float y, float r, int kind) {
        int[] colors = {0xFFE7B34D,0xFFD94E5C,0xFFB0D6E8,0xFF9A63CE,0xFF4BB6D5,0xFFE87E3E};
        p.setColor(Color.rgb(Color.red(colors[kind % colors.length]), Color.green(colors[kind % colors.length]), Color.blue(colors[kind % colors.length])));
        c.drawOval(new RectF(x - r, y - r * .62f, x + r, y + r * .62f), p);
        PathTri(c, x - r, y, x - r * 1.65f, y - r * .55f, x - r * 1.65f, y + r * .55f);
        p.setColor(Color.WHITE); c.drawCircle(x + r * .48f, y - r * .12f, r * .13f, p);
        p.setColor(Color.BLACK); c.drawCircle(x + r * .50f, y - r * .12f, r * .07f, p);
    }

    private void PathTri(Canvas c, float x1,float y1,float x2,float y2,float x3,float y3){
        android.graphics.Path path = new android.graphics.Path(); path.moveTo(x1,y1); path.lineTo(x2,y2); path.lineTo(x3,y3); path.close(); c.drawPath(path,p);
    }

    private void text(Canvas c, String v, float x, float y, float size, int color, boolean bold) { text(c,v,x,y,size,color,bold,false); }
    private void text(Canvas c, String v, float x, float y, float size, int color, boolean bold, boolean center) {
        p.setStyle(Paint.Style.FILL); p.setColor(color); p.setTextSize(size); p.setTypeface(bold ? android.graphics.Typeface.DEFAULT_BOLD : android.graphics.Typeface.DEFAULT);
        p.setTextAlign(center ? Paint.Align.CENTER : Paint.Align.LEFT); c.drawText(v,x,y,p);
    }

    @Override public boolean onTouchEvent(MotionEvent e) {
        if (e.getAction() != MotionEvent.ACTION_UP) return true;
        float x = e.getX(), y = e.getY(), w = getWidth(), h = getHeight();
        if (y > h - s(78)) {
            if (x < s(120)) screen = 0; else if (x < s(235)) screen = 1; else screen = 2;
            selectedFish = -1; return true;
        }
        if (screen == 0) {
            if (y > h * .28f && y < h * .40f) {
                if (x < s(165)) spot = 0; else if (x < s(315)) spot = 1; else spot = 2;
                return true;
            }
            if (selectedFish >= 0) {
                if (y > h * .82f && x > s(180) && x < s(335)) {
                    coins += VALUES[fish.get(selectedFish).kind];
                    selectedFish = -1;
                    return true;
                }
                if (y > h * .82f && x < s(175)) {
                    selectedFish = -1; return true;
                }
            } else if (y > h - s(120)) {
                int idx = random.nextInt(fish.size());
                selectedFish = idx;
                int k = fish.get(idx).kind; codex[k] = true; caught++;
                return true;
            } else if (y > h * .45f) {
                float fx = x / w;
                int nearest = 0; float dist = 99;
                for (int i = 0; i < fish.size(); i++) {
                    float d = Math.abs(fx - fish.get(i).x); if (d < dist) { dist = d; nearest = i; }
                }
                if (dist < .16f) { selectedFish = nearest; codex[fish.get(nearest).kind] = true; caught++; }
            }
        }
        return true;
    }

    private static class Fish {
        int kind; float x, speed;
        Fish(int kind, float x, float speed) { this.kind=kind; this.x=x; this.speed=speed; }
    }
}

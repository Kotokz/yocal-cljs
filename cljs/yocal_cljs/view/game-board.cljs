(ns yocal-cljs.view.game-board
  (:require [clojure.string :as str]
            [re-frame.core :as re-frame]
            [re-com.util :refer [item-for-id]]
            [re-com.core :refer [h-box gap label button v-box title]]))

;; ------------------ game board ---------------------------------------------
(defn select-values
  [map ks]
  (remove nil? (reduce #(conj %1 (map %2)) [] ks)))
(defn seq-map [dices] (frequencies (select-values :val dices)))

; can be 3/1/1 (3 of a kind )  or 3/2 (full house)
(defn check3
  [seqs]
  (let [ord (sort (vals seqs))]
    (cond
      (clojure.set/subset? #{2 3} (set ord)) (str "Full House!")
      :else (str "Three Of A Kind!"))))

; find the highest number with 2 occurrence
(defn check2
  [seqs]
  (let [[die crt] (first (sort > (filter #(> (second %) 1) seqs)))]
    (if (nil? die)
      (let [num (apply max (keys seqs))] (str "you have one of " num))
      (str "you have " crt " of " die))))

; check if small run or stright can be matched, otherwise check 2 of a kind
(defn check1
  [seqs]
  (let [ord (sort (keys seqs))]
    (cond
      (clojure.set/subset? #{1 2 3 4 5} (set ord)) (str "Long Run! You have 1 to 5")
      (clojure.set/subset? #{2 3 4 5 6} (set ord)) (str "Long Run! You have 2 to 6")
      (clojure.set/subset? #{1 2 3 4} (set ord)) (str "Small Run! from 1 to 4")
      (clojure.set/subset? #{2 3 4 5} (set ord)) (str "Small Run! from 2 to 5")
      (clojure.set/subset? #{3 4 5 6} (set ord)) (str "Small Run! from 3 to 6")
      :else (check2 seqs))))

(defn score
  [dices]
  (let [seqs (seq-map dices)
        [die occ] (first (sort-by last > seqs))]
    (case occ
      5 (str "Five Of A Kind! you have " occ " of " die)
      4 (str "Four Of A Kind! you have " occ " of " die)
      3 (check3 seqs)
      (check1 seqs))))

(defn roll-all
  [dices]
  (re-frame/dispatch [:set-roll-count])
  (let [nothelds (->> dices (filter (complement :isHeld)))]
    (if (-> dices count pos?)
      (doseq [ds nothelds]
        (re-frame/dispatch [:roll-dice (:id ds)]))
      (dotimes [id 5]
        (re-frame/dispatch [:add-dice id])))))

(defn dice
  [val]
  (case val
    1 [:div.first-face [:span.pip]]
    2 [:div.second-face [:span.pip]
       [:span.pip]]
    3 [:div.third-face [:span.pip]
       [:span.pip]
       [:span.pip]]
    4 [:div.fourth-face [:div.column [:span.pip]
                         [:span.pip]]
       [:div.column [:span.pip]
        [:span.pip]]]
    5 [:div.fifth-face [:div.column [:span.pip]
                        [:span.pip]]
       [:div.column [:span.pip]]
       [:div.column [:span.pip]
        [:span.pip]]]
    6 [:div.sixth-face [:div.column [:span.pip]
                        [:span.pip]
                        [:span.pip]]
       [:div.column [:span.pip]
        [:span.pip]
        [:span.pip]]]))

(defn dices [{:keys [id val isHeld]}]
  [:div.die-container {:on-click #(re-frame/dispatch [:hold-dice id])
                       :class    (if isHeld "isHeld")}
   [dice val]])

(defn game-board []
  [h-box
   :gap "2em"
   :children [(let [jwt (:jwt @(re-frame/subscribe [:user]))]
                (if (str/blank? jwt) (re-frame/dispatch [:set-active-panel :login-panel])))
              (let [ds (vals @(re-frame/subscribe [:game-dices]))]
                (if (-> ds count pos?)
                  [h-box
                   :children [(for [d ds]
                                ^{:key (:id d)} [dices d])
                              [:button {:on-click #(roll-all ds)} "Roll Again"]
                              [:button {:on-click #(re-frame/dispatch [:set-score-string (score ds)])} "Score"]]]
                  [:button {:on-click #(roll-all ds)} "Start game"]))
              (let [counts (re-frame/subscribe [:game-rolls])
                    score (re-frame/subscribe [:game-score])
                    score-string (re-frame/subscribe [:game-score-string])]
                [v-box
                 :gap "2em"
                 :children [
                            [h-box
                             :gap "2em"
                             :children [[:div (str "Roll counts = " @counts)]
                                        [:div (str "Game Score = " @score)]]]
                            [title :level :level2 :label @score-string]]])]])

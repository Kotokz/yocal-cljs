(ns yocal-cljs.db
  (:require [cljs.reader]
            [schema.core :as s :include-macros true]))

;; -- Schema -----------------------------------------------------------------
;;
;; This is a Prismatic Schema which documents the structure of app-db
;; See: https://github.com/Prismatic/schema
;;
(def schema {
              ;; Project name which should be string
              :name s/Str

              ;; a sorted-map is used to hold the dices.
              :game-dices (s/both PersistentTreeMap ;; ensure sorted-map, not just map
                            ;; each todo is keyed by its integer :id value
                            {s/Int {:id s/Int :val s/Int :isHeld s/Bool}})

              ;; Keep the score for Dice Game
              :game-score s/Int

              ;; Count the total rolls
              :game-rolls s/Int

              ;; The String for the game score pattern
              :game-score-string s/Str
              ;;
              :active-panel (s/enum :home-panel :about-panel :login-panel "")

              :user {:username s/Str :authenticated s/Bool :jwt s/Str}})

;; -- Default app-db Value  ---------------------------------------------------
;;
;; When the application first starts, this will be the value put in app-db
;; Look in core.cljs for  "(dispatch-sync [:initialise-db])"
;;
(def default-db
  {:name "YoCal"
   :game-dices (sorted-map)
   :game-score 0
   :game-rolls 0
   :game-score-string "Please start the game"
   :active-panel ""
   :user {:username ""
          :authenticated false
          :jwt ""}})

;; -- Local Storage  ----------------------------------------------------------
;;
(def lsk "yocal") ;; localstore key

(defn ls->yocal
  "Read in yocal from LS, and process into a map we can merge into app-db."
  []
  (some->> (.getItem js/localStorage lsk)
    (cljs.reader/read-string) ;; stored as an EDN map.
    (into (sorted-map)) ;; map -> sorted-map
    (hash-map :game-dices))) ;; access via the :todos key

(defn todos->ls!
  "Puts todos into localStorage"
  [game-dices]
  (.setItem js/localStorage lsk (str game-dices))) ;; sorted-map writen as an EDN map

(ns yocal-cljs.handlers
    (:require [re-frame.core :refer [register-handler path trim-v after dispatch debug log-ex]]
              [yocal-cljs.db :as db]
              [schema.core   :as s]))

;; -- Middleware --------------------------------------------------------------
;;

(defn check-and-throw
      "throw an exception if db doesn't match the schema."
      [a-schema data]
      (if-let [problems  (s/check a-schema data)]
         (throw (js/Error. (str "schema check failed: " problems)))))

;; after an event handler has run, this middleware can check that
;; it the value in app-db still correctly matches the schema.
(def check-schema-mw (after (partial check-and-throw db/schema)))
(def dice-middleware [check-schema-mw ;; after ever event handler make sure the schema is still valid
                      (path :game-dices)
                      debug
                      trim-v          ;; remove event id from event vec
                      ])

;; -- Helpers -----------------------------------------------------------------


(register-handler
 :initialize-db
 check-schema-mw
 (fn  [_ _]
   db/default-db))

(register-handler
 :set-active-panel
 debug
 (fn [db [_ active-panel]]
   (assoc db :active-panel active-panel)))

(register-handler
  :set-roll-count
  [(path :game-rolls)]
  (fn [game-rolls _]
    (inc game-rolls)))

(register-handler
  :add-dice
  dice-middleware
  (fn [game-dices [id]]
    (let [val (inc (rand-int 6))]
      (assoc game-dices id {:id id :val val :isHeld false}))))

(register-handler
  :roll-dice
  dice-middleware
  (fn [game-dices [id]]
    (let [val (inc (rand-int 6))]
        (assoc-in game-dices [id :val] val))))

(register-handler
  :hold-dice
  dice-middleware
  (fn [game-dices [id]]
    (update-in game-dices [id :isHeld] not)))

(register-handler
  :set-score-string
  [(path :game-score-string) debug trim-v]
  (fn [game-score-string [text]]
    text))

(ns yocal-cljs.subs
    (:require-macros [reagent.ratom :refer [reaction]])
    (:require [re-frame.core :as re-frame]))

;; -- Helpers -----------------------------------------------------------------



;; -- Subscription handlers and registration  ---------------------------------
(re-frame/register-sub
 :name
 (fn [db]
   (reaction (:name @db))))

(re-frame/register-sub
 :active-panel
 (fn [db _]
   (reaction (:active-panel @db))))

(re-frame/register-sub
 :game-dices
 (fn [db _]
  (reaction (:game-dices @db))))

(re-frame/register-sub
 :game-rolls
 (fn [db _]
  (reaction (:game-rolls @db))))

(re-frame/register-sub
 :game-score
 (fn [db _]
  (reaction (:game-score @db))))

(re-frame/register-sub
 :game-score-string
 (fn [db _]
  (reaction (:game-score-string @db))))

(re-frame/register-sub
  :user
  (fn [db _]
    (reaction (:user @db))))

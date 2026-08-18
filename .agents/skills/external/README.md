# external skills

本目录存放从外部仓库安装的第三方技能（如 obra/superpowers、mattpocock/skills），
通过 `skills-lock.json` 锁定来源与哈希。

安装方式：使用项目约定的技能安装工具，将技能安装到 `.agents/skills/external/<name>/`，
并由 `script/sync-skills-symlinks.sh` 同步到 IDE 目录。

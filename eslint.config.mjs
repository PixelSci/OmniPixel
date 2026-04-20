import antfu from '@antfu/eslint-config'

export default antfu({
    pnpm: true,
    yaml: false,

    ignores: ['apps/omni-pixel/src/components/ui/*'],

    stylistic: {
        indent: 4,
    },
})
